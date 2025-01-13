package tesla

import "sort"

func GetSortedVehicles(products []ProductInfo) []ProductInfo {
	// filter by "device_type": "vehicle"
	filtered := []ProductInfo{}
	for _, product := range products {
		if product.DeviceType == "vehicle" {
			filtered = append(filtered, product)
		}
	}

	// sort by "access_type": "OWNER"
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].AccessType == "OWNER" && filtered[j].AccessType != "OWNER"
	})

	return filtered
}
