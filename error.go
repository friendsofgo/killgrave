package killgrave

type invalidDirectoryError string

func (e invalidDirectoryError) Error() string { return string(e) }

type malformattedImposterError string

func (e malformattedImposterError) Error() string { return string(e) }
