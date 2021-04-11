package gettext

type PoEntries map[string]string

func (o PoEntries) Add(s string, entry string) (bool, *string) {
	if old, found := o[s]; found {
		return false, &old
	}
	o[s] = entry
	return true, nil
}
