package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type tempUser struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
	Name      string
}

type Users struct {
	List []tempUser `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	var (
		limit      int
		offset     int
		query      string
		orderField string
		orderBy    int
		err        error
	)
	authHeader := r.Header.Get("AccessToken")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No Authorization header", http.StatusUnauthorized)
	}
	if r.URL.Path != "/" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	strLimit := r.FormValue("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			fmt.Println("limit convert to int error", err)
		}
	}

	strOffset := r.URL.Query().Get("offset")
	if strOffset != "" {
		offset, err = strconv.Atoi(strOffset)
		if err != nil {
			fmt.Println("offset convert to int error", err)
		}
	}

	query = r.URL.Query().Get("query")

	orderField = r.URL.Query().Get("order_field")
	if orderField != "Name" && orderField != "Id" && orderField != "Age" {
		errStr, err := json.Marshal(SearchErrorResponse{"ErrorBadOrderField"})
		if err != nil {
			http.Error(w, "Error with marshaling", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errStr)
		return
	}

	strOrderBy := r.URL.Query().Get("order_by")
	if strOrderBy != "" {
		orderBy, err = strconv.Atoi(strOrderBy)
		if err != nil {
			fmt.Println("orderby convert to int error", err)
		}
	}

	// Открываем файл записывает адрес файла?
	file, err := os.Open("dataset.xml")
	if err != nil {
		fmt.Println("Error openin file:", err)
		return
	}
	defer file.Close()

	// Читаем файл ReadAll возвращает []byte
	xmlData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	var users Users
	err = xml.Unmarshal(xmlData, &users)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return
	}
	var rsltUser []User
	for _, userNode := range users.List {
		userNode.Name = userNode.FirstName + " " + userNode.LastName
		if strings.Contains(strings.ToLower(userNode.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(userNode.About), strings.ToLower(query)) {

			rsltUser = append(rsltUser, User{
				Id:     userNode.Id,
				Name:   userNode.Name,
				Age:    userNode.Age,
				About:  userNode.About,
				Gender: userNode.Gender,
			})
		}
	}

	sort.Slice(rsltUser, func(i, j int) bool {
		switch orderBy {
		case OrderByAsc:
			switch orderField {
			case "Id":
				return rsltUser[i].Id < rsltUser[j].Id
			case "Age":
				return rsltUser[i].Age < rsltUser[j].Age
			case "Name":
				return rsltUser[i].Name < rsltUser[j].Name
			default:
				return rsltUser[i].Name < rsltUser[j].Name
			}
		case OrderByAsIs:
			return false
		case OrderByDesc:
			switch orderField {
			case "Id":
				return rsltUser[i].Id > rsltUser[j].Id
			case "Age":
				return rsltUser[i].Age > rsltUser[j].Age
			case "Name":
				return rsltUser[i].Name > rsltUser[j].Name
			default:
				return rsltUser[i].Name > rsltUser[j].Name
			}
		default:
			return false
		}
	})

	lenList := len(rsltUser)
	lastUserIndex := limit + offset
	if lastUserIndex >= lenList {
		lastUserIndex = lenList - 1
	}
	w.Header().Set("Content-Type", "application/json")
	if lenList > 0 && offset < lenList {
		b, err := json.Marshal(rsltUser[offset:lastUserIndex])
		if err != nil {
			http.Error(w, "Error with marshaling", http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

// Как раз функция которая в client.go SearchClient тоже самое
// func checkResult(sr SearchRequest) (*SearchResponse, error) {
// 	url := fmt.Sprintf("127.0.0.1/?limit=%v&offset=%v&query=%v&orderBy=%v&orderField=%v", sr.Limit, sr.Offset, sr.Query, sr.OrderBy, sr.OrderField)
// 	resp, err := http.Get(url)

// 	if err != nil {
// 		return nil, err
// 	}
// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result := &SearchResponse{}

//		err = json.Unmarshal(data, result.Users)
//		if err != nil {
//			return nil, err
//		}
//		return result, nil
//	}

type TestCase struct {
	URL         string
	AccessToken string
	SR          SearchRequest
}

func TestSearchServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	cases := []TestCase{
		{
			URL:         ts.URL,
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      5,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL,
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      -1,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL,
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      85,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL,
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      5,
				Offset:     -2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL,
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      5,
				Offset:     2,
				Query:      "Sungfdfg1лапt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL,
			AccessToken: "",
			SR: SearchRequest{
				Limit:      5,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL + "/123/",
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      5,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         "127.0.0.2",
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      5,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
		},
		{
			URL:         ts.URL,
			AccessToken: "1",
			SR: SearchRequest{
				Limit:      5,
				Offset:     2,
				Query:      "Sunt",
				OrderField: "IDA",
				OrderBy:    OrderByAsc,
			},
		},
	}

	for numCase, item := range cases {
		sc := &SearchClient{
			URL:         item.URL,
			AccessToken: item.AccessToken,
		}
		result, err := sc.FindUsers(item.SR)

		if err != nil {
			fmt.Printf("[%d] error: %#v", numCase, err)
			continue
		}
		if result != nil {
			_ = ""
			// fmt.Println(result)
		}

	}
	ts.Close()
}

// func main() {
// 	http.HandleFunc("/", SearchServer)

// 	fmt.Println("starting server at :8080")
// 	http.ListenAndServe(":8080", nil)
// }
