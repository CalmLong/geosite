package main

func output(action int, originalMap map[string]struct{}) {
	for k := range originalMap {
		if v, ok := cover(k, action); ok {
			delete(originalMap, k)
			originalMap[v] = struct{}{}
		}
	}
}

func outputs(action int, originalMaps ...map[string]struct{}) {
	for i := 0; i < len(originalMaps); i++ {
		output(action, originalMaps[i])
	}
}
