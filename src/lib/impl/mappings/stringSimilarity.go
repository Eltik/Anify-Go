package mappings

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Rating struct {
	Target string
	Rating float64
}

type StringResult struct {
	Ratings        []Rating
	BestMatch      Rating
	BestMatchIndex int
}

type SimilarityResult struct {
	Same  bool
	Value float64
}

func Similarity(externalTitle, title string, titleArray []string) SimilarityResult {
	if title == "" {
		title = ""
	}
	simi := CompareTwoStrings(Clean(strings.ToLower(title)), strings.ToLower(externalTitle))
	for _, el := range titleArray {
		if el != "" {
			tempSimi := CompareTwoStrings(strings.ToLower(title), strings.ToLower(el))
			if tempSimi > simi {
				simi = tempSimi
			}
		}
	}

	found := simi > 0.6
	return SimilarityResult{Same: found, Value: simi}
}

func FindBestMatch(mainString string, targetStrings []string) StringResult {
	ratings := []Rating{}
	bestMatchIndex := 0

	for i, targetString := range targetStrings {
		currentRating := CompareTwoStrings(mainString, targetString)
		ratings = append(ratings, Rating{Target: targetString, Rating: currentRating})
		if currentRating > ratings[bestMatchIndex].Rating {
			bestMatchIndex = i
		}
	}

	bestMatch := ratings[bestMatchIndex]
	return StringResult{Ratings: ratings, BestMatch: bestMatch, BestMatchIndex: bestMatchIndex}
}

// FindBestMatchArray finds the best match from multiple main strings and target strings
func FindBestMatchArray(mainStrings, targetStrings []string) StringResult {
	mainStringResults := []StringResult{}

	for _, mainString := range mainStrings {
		ratings := []Rating{}
		bestMatchIndex := 0

		for i, targetString := range targetStrings {
			currentRating := CompareTwoStrings(mainString, targetString)
			ratings = append(ratings, Rating{Target: targetString, Rating: currentRating})

			if currentRating > ratings[bestMatchIndex].Rating {
				bestMatchIndex = i
			}
		}

		mainStringResults = append(mainStringResults, StringResult{
			Ratings:        ratings,
			BestMatch:      ratings[bestMatchIndex],
			BestMatchIndex: bestMatchIndex,
		})
	}

	overallBestMatchIndex := 0
	for i := 1; i < len(mainStringResults); i++ {
		if mainStringResults[i].BestMatch.Rating > mainStringResults[overallBestMatchIndex].BestMatch.Rating {
			overallBestMatchIndex = i
		}
	}

	return mainStringResults[overallBestMatchIndex]
}

// FindBestMatch2DArray finds the best match from a 2D target string array.
func FindBestMatch2DArray(mainStrings []string, targetStrings [][]string) StringResult {
	overallBestMatch := StringResult{
		Ratings:        []Rating{},
		BestMatch:      Rating{Target: "", Rating: 0},
		BestMatchIndex: 0,
	}

	for _, mainString := range mainStrings {
		for targetArrayIndex, targetArray := range targetStrings {
			ratings := []Rating{}

			for _, targetString := range targetArray {
				currentRating := CompareTwoStrings(Clean(strings.ToLower(strings.TrimSpace(mainString))), Clean(strings.ToLower(strings.TrimSpace(targetString))))
				ratings = append(ratings, Rating{Target: targetString, Rating: currentRating})
			}

			bestMatchIndex := 0
			for i, r := range ratings {
				if r.Rating > ratings[bestMatchIndex].Rating {
					bestMatchIndex = i
				}
			}

			if ratings[bestMatchIndex].Rating > overallBestMatch.BestMatch.Rating {
				overallBestMatch = StringResult{
					Ratings:        ratings,
					BestMatch:      ratings[bestMatchIndex],
					BestMatchIndex: targetArrayIndex,
				}
			}
		}
	}

	return overallBestMatch
}

// CompareTwoStrings calculates similarity between two strings using bigrams.
func CompareTwoStrings(first, second string) float64 {
	first = strings.ReplaceAll(first, " ", "")
	second = strings.ReplaceAll(second, " ", "")

	if first == second {
		return 1 // identical or empty
	}
	if len(first) < 2 || len(second) < 2 {
		return 0 // if either is a 0-letter or 1-letter string
	}

	firstBigrams := make(map[string]int)
	for i := 0; i < len(first)-1; i++ {
		bigram := first[i : i+2]
		firstBigrams[bigram]++
	}

	intersectionSize := 0
	for i := 0; i < len(second)-1; i++ {
		bigram := second[i : i+2]
		if count, exists := firstBigrams[bigram]; exists && count > 0 {
			firstBigrams[bigram]--
			intersectionSize++
		}
	}

	return 2.0 * float64(intersectionSize) / float64(len(first)+len(second)-2)
}

// Clean prepares a string for comparison by removing unnecessary characters and terms.
func Clean(title string) string {
	title = RemoveSpecialChars(title)
	title = TransformSpecificVariations(title)
	return title
}

// RemoveSpecialChars removes special characters from a string.
func RemoveSpecialChars(title string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9!@#$%^&*()\-= ]`)
	title = re.ReplaceAllString(title, " ")
	re2 := regexp.MustCompile(`[^A-Za-z0-9\-= ]`)
	title = re2.ReplaceAllString(title, "")
	title = strings.ReplaceAll(title, "  ", " ")
	return title
}

// TransformSpecificVariations standardizes specific variations in the text.
func TransformSpecificVariations(title string) string {
	title = strings.ReplaceAll(title, "yuu", "yu")
	title = strings.ReplaceAll(title, " ou", " oh")
	return title
}

func Slugify(args ...interface{}) string {
	// Define the default replacements as regular expressions
	replacements := map[string]string{
		`[aàáâãäåāăąǻάαа]`:   "a",
		`[bбḃ]`:              "b",
		`[cçćĉċčћ]`:          "c",
		`[dðďđδдђḋ]`:         "d",
		`[eèéêëēĕėęěέεеэѐё]`: "e",
		`[fƒφфḟ]`:            "f",
		`[gĝğġģγгѓґ]`:        "g",
		`[hĥħ]`:              "h",
		`[iìíîïĩīĭįıΐήίηιϊийіїѝ]`: "i",
		`[jĵј]`:            "j",
		`[kķĸκкќ]`:         "k",
		`[lĺļľŀłλл]`:       "l",
		`[mμмṁ]`:           "m",
		`[nñńņňŉŋνн]`:      "n",
		`[oòóôõöōŏőοωόώо]`: "o",
		`[pπпṗ]`:           "p",
		`q`:                "q",
		`[rŕŗřρр]`:         "r",
		`[sśŝşšſșςσсṡ]`:    "s",
		`[tţťŧțτтṫ]`:       "t",
		`[uùúûüũūŭůűųуў]`:  "u",
		`[vβв]`:            "v",
		`[wŵẁẃẅ]`:          "w",
		`[xξ]`:             "x",
		`[yýÿŷΰυϋύыỳ]`:     "y",
		`[zźżžζз]`:         "z",
		`[æǽ]`:             "ae",
		`[χч]`:             "ch",
		`[ѕџ]`:             "dz",
		`ﬁ`:                "fi",
		`ﬂ`:                "fl",
		`я`:                "ia",
		`[ъє]`:             "ie",
		`ĳ`:                "ij",
		`ю`:                "iu",
		`х`:                "kh",
		`љ`:                "lj",
		`њ`:                "nj",
		`[øœǿ]`:            "oe",
		`ψ`:                "ps",
		`ш`:                "sh",
		`щ`:                "shch",
		`ß`:                "ss",
		`[þθ]`:             "th",
		`ц`:                "ts",
		`ж`:                "zh",
		`[\\u0009-\\u000D\\u001C-\\u001F\\u0020\\u002D\\u0085\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\u058A\\u05BE\\u1400\\u1806\\u2010-\\u2015\\u2E17\\u2E1A\\u2E3A\\u2E3B\\u2E40\\u301C\\u3030\\u30A0\\uFE31\\uFE32\\uFE58\\uFE63\\uFF0D]`: "-",
	}

	// Join the arguments into a single string with spaces
	value := ""
	for _, arg := range args {
		value += fmt.Sprintf("%v ", arg)
	}
	value = strings.TrimSpace(value)

	// Apply each replacement
	for pattern, replacement := range replacements {
		re := regexp.MustCompile(pattern)
		value = re.ReplaceAllString(value, replacement)
	}

	// Normalize the string by removing accents, etc.
	value = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, value)

	// Lowercase, replace non-alphanumeric characters with dashes, and condense multiple spaces to a single dash
	value = strings.ToLower(value)
	value = regexp.MustCompile(`[^a-z0-9 ]+`).ReplaceAllString(value, "-")
	value = regexp.MustCompile(`\s+`).ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}
