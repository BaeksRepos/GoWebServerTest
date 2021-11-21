package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 구조체 형식으로 핸들러 함수 생성
// 인스턴스 생성 후 Handle에 등록
type fooHandler struct{}

// func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello, Foo!")
// }

// reqeust json to parsing data
type User struct {
	FirstName string    `json:"first_name"` //언어테이션
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// User 인스턴스 생성
	user := new(User)

	// Request의 Body는 Reader형태이며, NewDecoder와 Decode를통해 User형태로 변환
	// Decode(인스턴스)를 하면 Body의 데이터를 해당 인스턴스에 채워 넣어줌
	// Body는 IO Reader와 Closer를 가진 ReadCloser를 implement함
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}

	user.CreatedAt = time.Now()

	// 특정 인터페이스를 json형태로 인코딩
	// 결과 값은 Byte형태의 배열, Error를 반환
	data, _ := json.Marshal(user)

	// response header에 데이터 형태를 추가
	w.Header().Add("content-type", "application/json")

	// 성공일 경우 상태 OK 값 출력
	w.WriteHeader(http.StatusOK)

	// Response에 Byte배열인 data를 문자열로 변환후 response에 입력
	fmt.Fprint(w, string(data))

}

// // Handler Func 생성 후 HandleFunc에 등록하는 방식
func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {

	// http를 이용한 URL을 정적으로 등록

	// 즉 끝에 URL + / 일 경우의 작업처리
	// HandleFunc: 웹 Http 작업 핸들러 등록 함수
	// 어떤 Request 핸들러를 사용할지 지정하는 라우팅 함수
	// ResponseWriter: Request에 대한 응답 값
	// http.Request: 클라이언트에서 요청한 정보
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")

	})

	http.HandleFunc("/bar", barHandler)

	http.Handle("/foo", &fooHandler{})

	// 지정된 소프트에 웹서버를 열고 클라이언트 Request를 받아들여 작업을 할당
	http.ListenAndServe(":3000", nil)

	// mux 라우터 함수를 이용한 등록
	// 라우터 함수: 경로에 따라 작업에 대한 분배를 해주는 역할
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")

	})

	mux.HandleFunc("/bar", barHandler)
	mux.Handle("/foo", &fooHandler{})

	http.ListenAndServe(":3000", mux)

	//http.ListenAndServe(":3000", test)

}
