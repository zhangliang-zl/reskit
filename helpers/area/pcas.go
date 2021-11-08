package area

// pacs[province][city][area][]streets
var pacs map[string]map[string]map[string][]string

func PACS() map[string]map[string]map[string][]string {
	return pacs
}

func Provinces() (provinces []string) {
	for province := range pacs {
		provinces = append(provinces, province)
	}

	return provinces
}

func Cities(p string) (cities []string) {
	if children, ok := pacs[p]; ok {
		for city := range children {
			cities = append(cities, city)
		}
	}

	return cities
}

func Areas(province, city string) (areas []string) {
	if _, ok := pacs[province]; !ok {
		return
	}

	children, ok := pacs[province][city]
	if !ok {
		return
	}

	for area := range children {
		areas = append(areas, area)
	}
	return
}

func Streets(province, city, area string) (streets []string) {
	if _, ok := pacs[province]; !ok {
		return
	}

	if _, ok := pacs[province][city]; !ok {
		return
	}

	streets, ok := pacs[province][city][area]
	if !ok {
		return
	}

	return
}
