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

func SimpleTestSubConfigsJSON() []byte {
	return []byte(`
		{
			"articles": {
				"entity_type": "articles",
				"bundle": "",
				"mapping": {
					"author": {
						"type": "relation",
						"name": "author"
					},
					"comments": {
						"type": "relation",
						"name": "comments"
					}
				}
			}, 
			"people": {
				"entity_type": "people",
				"bundle": "",
				"mapping": {}
			},
			"comments": {
				"entity_type": "people",
				"bundle": "",
				"mapping": {
					"author": {
						"type": "relation",
						"name": "author"
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
