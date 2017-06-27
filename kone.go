package kons

import (
	"math"
	"regexp"
	"sync"
	"sync/atomic"
)

type Kone struct {
	root  *Kone
	dmu   sync.RWMutex
	data  interface{}
	kmu   sync.RWMutex
	kones map[string]*Kone
}

func NewKone(root *Kone) *Kone {
	return &Kone{
		root: root,
	}
}

func (k *Kone) GetData() interface{} {
	k.dmu.RLock()
	defer k.dmu.RUnlock()

	return k.data
}

func (k *Kone) SetData(data interface{}) {
	k.dmu.Lock()
	defer k.dmu.Unlock()

	k.data = data
}

func (k *Kone) GetRoot() *Kone {
	if k.root != nil {
		return k.root
	}
	return k
}

func (k *Kone) GetKone(key string) *Kone {
	k.kmu.RLock()
	defer k.kmu.RUnlock()

	if kone, ok := k.kones[key]; ok {
		return kone
	}
	return nil
}

func (k *Kone) SetKone(key string, kone *Kone) {
	k.kmu.Lock()
	defer k.kmu.Unlock()

	if k.kones == nil {
		k.kones = make(map[string]*Kone)
	}
	k.kones[key] = kone
}

func (k *Kone) GetKones() map[string]*Kone {
	k.kmu.RLock()
	defer k.kmu.RUnlock()

	return k.kones
}

func (k *Kone) SetKones(kones map[string]*Kone) {
	k.kmu.Lock()
	defer k.kmu.Unlock()

	k.kones = kones
}

func (k *Kone) Upsert(settings *UpsertSettings, data interface{}, path ...string) (int, error) {
	if settings == nil {
		settings = DefaultUpsertSettings()
	}
	return k.upsert(settings, data, path)
}

func (k *Kone) Find(settings *FindSettings, regPath ...string) ([]*Kone, error) {
	regs, err := PathToRegs(regPath...)
	if err != nil {
		return nil, err
	}
	return k.FindReg(settings, regs...)
}

func (k *Kone) FindReg(settings *FindSettings, regs ...*regexp.Regexp) ([]*Kone, error) {
	if settings == nil {
		settings = DefaultFindSettings()
	}
	if settings.Limit == 0 {
		settings.Limit = math.MaxInt32
	}

	var n int32
	var kones []*Kone
	for kone := range k.findReg(settings, regs, &n) {
		kones = append(kones, kone)
	}

	if len(kones) > int(settings.Limit) {
		kones = kones[:settings.Limit]
	}
	return kones, nil
}

func (k *Kone) FindOne(settings *FindSettings, regPath ...string) (*Kone, error) {
	regs, err := PathToRegs(regPath...)
	if err != nil {
		return nil, err
	}
	return k.FindOneReg(settings, regs...)
}

func (k *Kone) FindOneReg(settings *FindSettings, regs ...*regexp.Regexp) (*Kone, error) {
	if settings == nil {
		settings = DefaultFindSettings()
	}
	settings.Limit = 1

	kones, err := k.FindReg(settings, regs...)
	if err != nil {
		return nil, err
	}
	if len(kones) == 0 {
		return nil, nil
	}
	return kones[0], nil
}

func (k *Kone) findReg(settings *FindSettings, regs []*regexp.Regexp, n *int32) <-chan *Kone {
	kones := make(chan *Kone, 1)
	var wg sync.WaitGroup
	if len(regs) == 0 {
		atomic.AddInt32(n, 1)
		kones <- k
	} else {
		reg := regs[0]
		for key, kone := range k.GetKones() {
			kone := kone
			if atomic.LoadInt32(n) > settings.Limit {
				break
			}
			if !reg.MatchString(key) {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				for kn := range kone.findReg(settings, regs[1:], n) {
					kones <- kn
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		close(kones)
	}()
	return kones
}

func (k *Kone) upsert(settings *UpsertSettings, data interface{}, path []string) (int, error) {
	if len(path) == 0 {
		k.SetData(data)
		return 1, nil
	}

	key := path[0]
	kone := k.GetKone(key)
	if kone == nil && settings.MakePath {
		kone = NewKone(k.GetRoot())
		k.SetKone(key, kone)
	}
	if kone != nil {
		return kone.upsert(settings, data, path[1:])
	}
	return 0, nil
}
