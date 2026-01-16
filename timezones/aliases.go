package timezones

// ManualAliases maps common abbreviations and nicknames to canonical city names
// These are manually curated for common use cases
var ManualAliases = map[string]string{
	// United States - Major Cities
	"nyc":            "new york",
	"ny":             "new york",
	"big apple":      "new york",
	"la":             "los angeles",
	"sf":             "san francisco",
	"chi":            "chicago",
	"philly":         "philadelphia",
	"dc":             "washington",
	"atl":            "atlanta",
	"hotlanta":       "atlanta",
	"bos":            "boston",
	"vegas":          "las vegas",
	"phx":            "phoenix",
	"pdx":            "portland",
	"sea":            "seattle",
	"det":            "detroit",
	"mia":            "miami",
	"dal":            "dallas",
	"hou":            "houston",
	"nola":           "new orleans",
	"the bay":        "san francisco",
	"silicon valley": "san jose",

	// Canada
	"to":    "toronto",
	"the 6": "toronto",
	"yvr":   "vancouver",
	"mtl":   "montreal",

	// United Kingdom
	"ldn":      "london",
	"the city": "london",

	// Europe
	"paname": "paris",
	"barca":  "barcelona",
	"bcn":    "barcelona",
	"mad":    "madrid",

	// Asia
	"hk":  "hong kong",
	"sg":  "singapore",
	"bkk": "bangkok",
	"del": "delhi",
	"bom": "mumbai",
	"blr": "bangalore",

	// Australia
	"syd":  "sydney",
	"mel":  "melbourne",
	"bris": "brisbane",

	// Middle East
	"dxb": "dubai",

	// Africa
	"jnb":     "johannesburg",
	"jo'burg": "johannesburg",
	"joburg":  "johannesburg",
	"cpt":     "cape town",
	"nbo":     "nairobi",
}
