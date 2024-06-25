package main

func ToObject(node Node) interface{} {
	switch nodeType := node.(type) {
	case *Object:
		obj := map[string]interface{}{}
		for k, v := range nodeType.Pairs {
			obj[k] = ToObject(v)
		}
		return obj
	case *Array:
		arr := make([]interface{}, len(nodeType.Elements))
		for i, el := range nodeType.Elements {
			arr[i] = ToObject(el)
		}
		return arr
	case *Value:
		return nodeType.Value
	}

	return nil
}
