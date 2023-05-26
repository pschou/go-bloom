package bloom

import "os"

func (f *Filter) SaveFile(file string) error {
	flt, err := os.Create(file)
	if err != nil {
		return err
	}
	defer flt.Close()
	return f.Save(flt)
}

func (f *Filter) LoadFile(file string) (*Filter, error) {
	flt, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer flt.Close()
	return Load(flt)
}
