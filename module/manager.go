package module

type Modules map[nameKey]Instance

func (m Modules) Value(key any) Instance {
	nameKey, ok := key.(nameKey)
	if !ok {
		return nil
	}

	return m[nameKey]
}
