package main

func output(action int, originalMap map[string]struct{}) {
	for k := range originalMap {
		if v, ok := toV2Ray(k, action); ok {
			delete(originalMap, k)
			originalMap[v] = struct{}{}
		}
	}
}