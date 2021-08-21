# Buddy Back-end Activity API Specification

0. Server Domain:Port

    http://localhost:3000 (추후 변경 예정)

<br>

1. Create - 활동 생성

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/create | manager |

    - Request
        - start: (string) 시작일, Unixtimestamp
        - end: (string) 종료일, Unixtimestamp
        - place: (string) 장소
        - type: (number) 활동 종류 (창립제: 0, 스터디: 1, 기타: 2)
        - description: (string) 활동 설명
        - participants: (Array&lt;string&gt;) 참여자 학번 목록
        - private: (bool) 해당 활동의 private 여부

    - Request Body example
        ```json
        {
            "start": "1628249722",
            "end": "1628250522",
            "place": "성곡도서관 2층 스터디실 2번방",
            "type": 1,
            "description": "2022년 1학기 2차 알고리즘 스터디",
            "participants" : [
                "20200201",
                "20191524",
                "20212282"
            ],
            "private": true
        }
        ```

    - Response
        - error: (string) 에러 메시지 (활동 생성 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string"
        }
        ```

    - Status Code
        - 200 OK: 활동 생성 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

2. Search - 활동 검색 (Landing Page)

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/activity/search | member |

    - Query Parameter
        - query: (string) 검색어

    - Query Parameter example
        ```json
        http://localhost:3000/api/v1/activity/search?query=2021
        ```

    - Response
        - data.activities: (Array&lt;string&gt;) 활동 검색 결과
        - error: (string) 에러 메시지 (활동 검색 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "activities": [
                    {
                        "id": "610d458b79e122ea1d150cd6",
                        "start": "1628249722",
                        "end": "1628249722",
                        "place": "공학관 209호",
                        "type": 0,
                        "description": "2021년 창립제",
                        "participants": [
                            {
                                "id": "20190000",
                                "name": "홍길동",
                                "department": "예술대학 도자기학과",
                                "email": "gildong@gmail.com",
                                "grade": 3
                            },
                            {
                                "id": "20175271",
                                "name": "김희동",
                                "department": "와플대학 아이스크림학과",
                                "email": "hdong@kookmin.ac.kr",
                                "grade": 1
                            }
                        ],
                        "private": false
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 500 Internal Server Error: 시스템 오류

3. Private Search - 활동 검색 (Back Office)

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/activity/private | manager |

    - Query Parameter
        - query: (string) 검색어

    - Query Parameter example
        ```json
        http://localhost:3000/api/v1/activity/private?query=2021
        ```

    - Response
        - data.activities: (Array&lt;string&gt;) 활동 검색 결과
        - error: (string) 에러 메시지 (활동 검색 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "activities": [
                    {
                        "id": "610d458b79e122ea1d150cd6",
                        "start": "1628249722",
                        "end": "1628249722",
                        "place": "성곡도서관 2층 스터디룸",
                        "type": 1,
                        "description": "2021년 2차 알고리즘 스터디",
                        "participants": [
                            {
                                "id": "20190000",
                                "password": "asdf232",
                                "name": "홍길동",
                                "department": "예술대학 도자기학과",
                                "phone": "010-1234-5678",
                                "email": "gildong@gmail.com",
                                "grade": 3,
                                "attendance": 0,
                                "approved": true,
                                "on_delete": false,
                                "created_at": "1628974315",
                                "updated_at": "1628974315",
                                "role": {
                                    "member_management": false,
                                    "activity_management": false,
                                    "fee_management": false
                                }
                            },
                            {
                                "id": "20175271",
                                "password": "sfssll292",
                                "name": "김희동",
                                "department": "와플대학 아이스크림학과",
                                "phone": "010-3131-3232",
                                "email": "hdong@kookmin.ac.kr",
                                "grade": 1,
                                "attendace": 1,
                                "approved": true,
                                "on_delete": false,
                                "created_at": "1629060720",
                                "updated_at": "1629060720",
                                "role": {
                                    "member_management": false,
                                    "activity_management": false,
                                    "fee_management": false
                                }
                            }
                        ],
                        "private": true
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 500 Internal Server Error: 시스템 오류

4. Update - 활동 정보 수정

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/activity/update | manager |

    - Request
        - id: (string) 수정할 활동 ID
        - update: (JSON) 수정할 활동 정보 (시작일, 종료일, 장소, 종류, 설명, 참여자 목록 중 0개 이상 택)

    - Request Body example
        ```json
        {
            "id": "610d458b79e122ea1d150cd6",
            "update": {
                "start": "1628249722",
                "end": "1628249722",
                "place": "cafe",
                "type": "study",
                "description": "Study End!",
                "participants": [
                    "20192019",
                    "20182018"
                ]
            }
        }
        ```

    - Response
        - error: (string) 에러 메시지 (활동 정보 갱신 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 활동 정보 갱신 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

5. Delete - 활동 삭제

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | DELETE | /api/v1/activity/delete | manager |

    - Request
        - id: (string) 삭제할 활동 ID

    - Request Body example
        ```json
        {
            "id": "610d458b79e122ea1d150cd6"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (활동 삭제 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "mongo: no documents in result"
        }
        ```

    - Status Code
        - 200 OK: 활동 삭제 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 잘못된 ID, 시스템 오류 등

6. Upload - 파일 업로드 (사진 포함)

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/upload | manager |

    - Query Parameter
        - id: (string) 활동 ID
    
    - Query Parameter example
        ```json
        http://localhost:3000/api/v1/activity/upload?id=6120347c7289f5bf7e22a7ad
        ```
    
    - Request form
        - file (최대 용량: 32MiB)
    
    - Request form example
        - [여기](https://github.com/kmu-kcc/buddy-backend/blob/develop/testutil/upload_test.html)를 참고하세요.
    
    - Response
        - error: (string) 에러 메시지 (파일 업로드 성공 시 empty)
    
    - Response Body example
        ```json
        {
            "error": "http:no such file"
        }
        ```

    - Status Code
        - 200 OK: 파일 업로드 성공
        - 500 Internal Server Error: 파일이 없는 경우, 저장 실패, 잘못된 ID, 시스템 오류 등

7. Download - 파일 다운로드

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/download | - |

    - Request
        - filename: (string) 다운받고자 하는 파일 이름

   - Request Body Example
    ```json
    {
        "filename": "motorcycle.svg"
    }
    ```

    - Response
        - error: (string) 에러 메시지 (400 Bad Request에만 해당)
        - file: 찾고자 하는 파일

    - Response Body example
        ```json
        "error": "404 page not found"
        ```

    - Status Code
        - 400 Bad Request: 요청 포맷/타입 오류
        - 404 Page Not Found: 찾고자 하는 파일이 없는 경우

8. Delete File - 파일 삭제

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/deletefile | manager |

    - Request
        - id: (string) 활동 ID
        - filename: (string) 삭제하고자 하는 파일명
    
    - Request Body example
        ```json
        {
            "id": "6120347c7289f5bf7e22a7ad",
            "filename": "motorcycle.svg"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (파일 삭제 성공 시 empty)
    
    - Response Body example
        ```json
        {
            "error": "http:no such file"
        }
        ```

    - Status Code
        - 200 OK: 파일 삭제 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 잘못된 ID, 시스템 오류 등