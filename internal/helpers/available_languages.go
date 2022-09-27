package helpers

// Languages is the available language list.
// var Languages, _ = GetLanguages().
var Languages = []string{
	"Unknown",
	"English",
	"French",
	"Spanish",
	"German",
	"Italian",
	"Danish",
	"Dutch",
	"Japanese",
	"Icelandic",
	"Chinese",
	"Russian",
	"Polish",
	"Vietnamese",
	"Swedish",
	"Norwegian",
	"Finnish",
	"Turkish",
	"Portuguese",
	"Flemish",
	"Greek",
	"Korean",
	"Hungarian",
	"Hebrew",
	"Lithuanian",
	"Czech",
	"Arabic",
	"Hindi",
	"Bulgarian",
	"Malayalam",
	"Ukrainian",
	"Slovak",
}

// GetLanguages pull languages from Sonarr source code and converts it to slice.
// using static slice to avoid github dependency.
// func GetLanguages() ([]string, error) {
// 	var languages []string

// 	resp, err := http.Get("https://raw.githubusercontent.com/Sonarr/Sonarr/develop/src/NzbDrone.Core/Languages/Language.cs")
// 	if err != nil {
// 		fmt.Println(err)

// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println(err)

// 		return nil, err
// 	}

// 	r := strings.NewReplacer(" ", "", "\n", "\"", ",", "\",", "{", "[", "}", "]")
// 	if err := json.Unmarshal([]byte(r.Replace(strings.Split(strings.Split(string(b), "return new List<Language>\n")[1], ";")[0])), &languages); err != nil {
// 		panic(err)
// 	}

// 	return languages, nil
// }

// GetLanguageID retrieve language ID of a given language.
func GetLanguageID(language string) int64 {
	languages := Languages
	for i, l := range languages {
		if l == language {
			return int64(i)
		}
	}

	return 0
}
