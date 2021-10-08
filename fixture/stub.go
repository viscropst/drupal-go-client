package fixture

func NodeBannerTestSubConfigsJSON() []byte {
	return []byte(`
		{
			"node--banner": {
				"entity_type": "node",
				"bundle": "banner",
				"mapping": {
					"field_banner_image": {
						"type": "file",
						"name": "banner_image"
					},
					"field_banner_link": {
						"type": "raw",
						"name": "link"
					}
				}
			}
		}
	`)
}

func NodeBannerTestNoMappingIgnoreSubConfigsJSON() []byte {
	return []byte(`
		{
			"node--banner": {
				"entity_type": "node",
				"bundle": "banner",
				"no_mapping_mode": "ignore",
				"mapping": {
					"field_banner_image": {
						"type": "file",
						"name": "banner_image"
					},
					"field_banner_link": {
						"type": "raw",
						"name": "link"
					}
				}
			}
		}
	`)
}
