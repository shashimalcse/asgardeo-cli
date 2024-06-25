package tui

type Item struct {
	key, title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) Key() string         { return i.key }
func (i Item) FilterValue() string { return i.title }

func NewItem(title, desc string) Item {
	return Item{title: title, desc: desc}
}

func NewItemWithKey(key, title, desc string) Item {
	return Item{key: key, title: title, desc: desc}
}
