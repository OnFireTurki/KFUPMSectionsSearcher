package LFS

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Options struct {
	Name, TeleID string
	WaitListMode bool
}
type Worker struct {
	CourseSub, CourseNume, Term, JSESSIONID string
	OptionSe                                Options
	CurrentData                             []CourseDetails
}

var (
	hed map[string]string = map[string]string{
		"Conncetion":                "keep-alive",
		"Upgrade-Insecure-Requests": "0",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "en",
		"Host":                      "banner9-registration.kfupm.edu.sa",
	}
)

func (W *Worker) activateJSID(id string) bool {
	req, _ := SendRequest("POST", "https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/term/search?mode=search", "term="+W.Term+"&studyPath=&studyPathText=&startDatepicker=&endDatepicker=", map[string]string{"Cookie": "JSESSIONID=" + id}, nil, false)
	return !strings.Contains(req.Body, "regAllowed\": false")
}

func (W *Worker) GetAndActivate() {
	req, _ := SendRequest("GET", "https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/term/termSelection?mode=search", "", hed, nil, false)
	for i := range req.Cookies {
		if i == "JSESSIONID" {
			W.JSESSIONID = req.Cookies[i]
			break
		}
	}
	W.activateJSID(W.JSESSIONID)
}

func (W *Worker) search() JsonRes {
	req, _ := SendRequest("GET", fmt.Sprintf("https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/searchResults/searchResults?txt_subject=%v&txt_courseNumber=%v&txt_term=%v&startDatepicker=&endDatepicker=&pageOffset=0&pageMaxSize=100&sortColumn=subjectDescription&sortDirection=asc&txt_campus=C", W.CourseSub, W.CourseNume, W.Term), "", map[string]string{"Cookie": "JSESSIONID=" + W.JSESSIONID}, nil, false)
	// txt_campus : F / For females stupid sections
	res := JsonRes{}
	json.Unmarshal([]byte(req.Body), &res)
	return res
}

func (W *Worker) filter(data JsonRes) {
	W.CurrentData = nil
	if W.OptionSe.Name != "" && W.OptionSe.WaitListMode == false {
		for i := 0; i < data.TotalCount; i++ {
			for j := 0; j < len(data.Data[i].Faculty); j++ {
				if strings.Contains(W.OptionSe.Name, data.Data[i].Faculty[j].DisplayName) {
					W.CurrentData = append(W.CurrentData, data.Data[i])
					break
				}
			}
		}
	} else { // I KNOW Missing statement wih the waitlistmode
		// EVERYTHING
		W.CurrentData = data.Data
	}

}

func (W *Worker) JIDThread() {
	// We need to update the jsession id every half minute
	for {
		W.GetAndActivate()
		time.Sleep(time.Second * 30)
	}
}

func (W *Worker) SingleSearch() []CourseDetails {
	W.GetAndActivate()
	W.filter(W.search())
	return W.CurrentData
}

// I decided to stop develop this project due to several problems in the Login and Registration part.
// But not fully stoping, I will move this project to another language,
// with focusing on the login/registersion more than the search/filtering part
// 6 / 26 / 2022
