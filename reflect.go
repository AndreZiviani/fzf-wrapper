package fzfwrapper

/*
func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}
func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}
*/

/*
// Converts fzf.Item to MyItem
func NewMyItem(item *fzf.Item) *MyItem {
	return &MyItem{
		text:        GetUnexportedField(reflect.ValueOf(item).Elem().FieldByName("text")).(util.Chars),
		transformed: GetUnexportedField(reflect.ValueOf(item).Elem().FieldByName("transformed")).(*[]fzf.Token),
		origText:    GetUnexportedField(reflect.ValueOf(item).Elem().FieldByName("origText")).(*[]byte),
	}
}

func NewMyItemList(items [fzfChunkSize]fzf.Item) []MyItem {
	ret := make([]MyItem, 0)
	for _, i := range items {
		tmp := NewMyItem(&i)
		ret = append(ret, *tmp)
	}

	return ret
}
*/
