# Buddy Back-end Member API Specification

0. Server Domain:Port

    http://localhost:3000 (추후 변경 예정)

<br>

* 로그인을 제외한 모든 API들은 헤더의 Authorization 필드의 access token을 통해 토큰 유효성과 API 접근 권한을 검사합니다.

* 만약 access token이 유효하지 않은 경우 401 Unauthorized, 해당 API에 대한 접근 권한이 없을 경우 403 Permission Denied 오류를 반환합니다.

* 발급받은 access token은 6시간 뒤에 자동으로 삭제됩니다. 토큰이 유효성 검증 실패 오류를 받으면 재로그인하도록 처리해주세요.

    - 401 Unauthorized: 토큰 인증 실패
    - 403 Permission Denied: 접근 권한 없음

<br>

1. SignIn - 회원 로그인 (로그인 성공 후 사용하는 모든 API들의 헤더의 Authorization 필드에 발급받은 access token을 넣어주세요 :) )

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/signin | member |

    - Request
        - id: (string) 학번
        - password: (string) 비밀번호

    - Request Body example
        ```json
        {
            "id": "20210000",
            "password": "asdf1234"
        }
        ```

    - Response
        - data.access_token: (string) 엑세스 토큰
        - data.expired_at: (string) 만료 시각 (Unix timestamp)
        - error: (string) 에러 메시지 (로그인 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHBpcmUiOjE2Mjk1NzMwOTMsImlkIjoiMjAyMTAwMDEifQ.LXE7k3kubkCMFJL7bNQYTWPDymCFtclKOBIAKHqDHfQ",
                "expired_at": "1629573093"
            },
            "error": "password mismatch"
        }
        ```

    - Status Code
        - 200 OK: 로그인 성공 (유효한 엑세스 토큰 발급)
        - 422 Unprocessable Entity: 엑세스 토큰 발급 실패
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: ID/PW 오류, 가입 미승인 상태, 시스템 오류 등

2. SignUp - 회원 가입 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/signup | - |

    - Request
        - id: (string) 학번
        - name: (string) 이름
        - department: (string) 소속 대학/학부
        - phone: (string) 전화번호
        - email: (string) 메일 주소
        - grade: (number) 학년 (1 이상의 정수)
        - attendance: (number) 재학 여부 (재학: 0, 휴학: 1: 졸업: 2)

    - Request Body example
        ```json
        {
            "id": "20210000",
            "name": "홍길동",
            "department": "조형대학 시각디자인학과",
            "phone": "010-1234-5678",
            "email": "gildong@gmail.com",
            "grade": 1,
            "attendance": 0
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회원 가입 신청 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "under review"
        }
        ```

    - Status Code
        - 200 OK: 회원 가입 신청 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 이미 회원인 경우, 가입 신청 처리 중인 경우, 시스템 오류 등

3. SignUps - 회원 가입 신청 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/signups | member manager |

    - Response
        - data.signups: (Array&lt;JSON&gt;) 회원 가입 신청자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "signups": [
                    {
                        "id": "20190000",
                        "password": "20190000",
                        "name": "김희동",
                        "department": "와플대학 팥빙수학과",
                        "phone": "010-1234-5678",
                        "email": "heedong@gmail.com",
                        "grade": 3,
                        "attendance": 0,
                        "approved": false,
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
                        "id": "20200299",
                        "password": "20200299",
                        "name": "이기철",
                        "department": "자연과학대학 물리학과",
                        "phone": "010-9876-5432",
                        "email": "lee@naver.com",
                        "grade": 2,
                        "attendance": 1,
                        "approved": false,
                        "on_delete": false,
                        "created_at": "1629060720",
                        "updated_at": "1629060720",
                        "role": {
                            "member_management": false,
                            "activity_management": false,
                            "fee_management": false
                        }
                    }
                ]
            },
            "error": "permission denied"
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 400 Bad Request: 요청 포맷/응답 오류
        - 500 Internal Server Error: 시스템 오류

4. Approve - 회원 가입 승인

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/member/approve | member manager |

    - Request
        - ids: (Array&lt;string&gt;) 회원 가입을 승인할 신청자들의 학번 List

    - Request Body example
        ```json
        {
            "ids": [
                "20210000",
                "20180020",
                "20170011"
            ]
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회원 가입 승인 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string"
        }
        ```

    - Status Code
        - 200 OK: 회원 가입 승인 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

5. Exit - 회원 탈퇴 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/member/exit | member |

    - Request
        - id: (string) 탈퇴 신청하는 회원의 학번

    - Request Body example
        ```json
        {
            "id": "20210000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (탈퇴 신청 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "already on delete"
        }
        ```

    - Status Code
        - 200 OK: 탈퇴 신청 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 삭제 요청 처리 중인 경우, 시스템 오류

6. Exits - 회원 탈퇴 신청 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/exits | member manager |

    - Response
        - data.exits: (Array&lt;JSON&gt;) 회원 탈퇴 신청자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "exits": [
                    {
                        "id": "20190000",
                        "password": "asdf1234",
                        "name": "김희동",
                        "department": "경상대학 스포츠레저학과",
                        "phone": "010-1234-5678",
                        "email": "heedong@gmail.com",
                        "grade": 3,
                        "attendance": 0,
                        "approved": true,
                        "on_delete": true,
                        "created_at": "1629060720",
                        "updated_at": "1629060930",
                        "role": {
                            "member_management": false,
                            "activity_management": false,
                            "fee_management": false
                        }
                    },
                    {
                        "id": "20200299",
                        "password": "abcdefg9",
                        "name": "이기철",
                        "department": "자연과학대학 물리학과",
                        "phone": "010-9876-5432",
                        "email": "lee@naver.com",
                        "grade": 2,
                        "attendance": 1,
                        "approved": true,
                        "on_delete": true,
                        "created_at": "1629080720",
                        "updated_at": "1629081720",
                        "role": {
                            "member_management": false,
                            "activity_management": true,
                            "fee_management": false
                        }
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 500 Internal Server Error: 시스템 오류

7. Delete - 회원 가입 거부 및 회원 탈퇴 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | DELETE | /api/v1/member/delete | member manager |

    - Request
        - ids: (Array&lt;string&gt;) 회원 가입 승인 거부/탈퇴 처리하는 신청자/회원들의 학번 List

    - Request Body example
        ```json
        {
            "ids": [
                "20210000",
                "20180020",
                "20170011"
            ]
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회원 가입 거부/탈퇴 처리 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string"
        }
        ```

    - Status Code
        - 200 OK: 회원 가입 거부/탈퇴 처리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

8. My - 내 정보

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/my | - |

    - Request
        - id: (string) 학번
        - password: (string) 비밀번호

    - Request Body example
        ```json
        {
            "id": "20210021",
            "password": "abcdffw112"
        }
        ```

    - Response
        - data.data: (JSON) 내 정보
        - error: (string) 에러 메시지 (회원 검색 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "data": {
                    "id": "202100021",
                    "password": "abcdffw112",
                    "name": "홍길동",
                    "department": "소프트웨어융합대학 소프트웨어학부",
                    "phone": "01012345678",
                    "email": "gildong@kookmin.ac.kr",
                    "grade": 1,
                    "attendance": 0,
                    "approved": true,
                    "on_delete": false,
                    "role": {
                        "member_management": false,
                        "activity_management": false,
                        "fee_management": false
                    }
                }
            },
            "error": "password mismatch"
        }
        ```

    - Status Code
        - 200 OK: 내 정보 가져오기 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 비밀 번호 오류, 시스템 오류 등

9. Search - 회원 검색

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/search | member |

    - Query Parameter
        - query: (string) 검색어

    - Query Parameter example
        ```json
        http://localhost:3000/api/v1/member/search?query=20190302
        ```

    - Response
        - data.members: (Array&lt;JSON&gt;) 회원 검색 결과
        - error: (string) 에러 메시지 (회원 검색 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "members": [
                    {
                        "id": "20210000",
                        "name": "홍길동",
                        "department": "소프트웨어융합대학 소프트웨어학부",
                        "email": "gildong@kookmin.ac.kr",
                        "grade": 1,
                        "role": {
                            "member_management": false,
                            "activity_management": false,
                            "fee_management": false
                        }
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 회원 검색 성공
        - 500 Internal Server Error: 시스템 오류

10. Update - 회원 정보 갱신

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/member/update | member |

    - Request
        - id: (string) 학번
        - update: (JSON) 갱신하고자 하는 회원 정보 (비밀번호, 소속 대학/학부, 전화번호, 이메일, 학년, 재학 여부 중 0개 이상 택)

    - Request Body example
        ```json
        {
            "id": "20210001",
            "update": {
                "password": "asdf1234",
                "department": "소프트웨어융합대학 소프트웨어학부",
                "phone": "010-1234-5678",
                "email": "gildong@yahoo.com",
                "grade": 2,
                "attendance": 0
            }
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회원 정보 갱신 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string"
        }
        ```

    - Status Code
        - 200 OK: 회원 정보 갱신 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

11. Active - 회원 가입 신청 활성 상태 확인

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/active | - |

    - Response
        - data.active: (boolean) 활성 여부
        - error: (string) 에러 메시지 (확인 성공 시 empty)

    - Response Body example
    ```json
    {
        "data": {
            "active": true
        },
        "error": "this MongoDB deployment does not support retryable writes. Please add retryWrites=false to your connection string"
    }
    ```

    - Status Code
        - 200 OK: 확인 성공
        - 500 Internal Server Error: 시스템 오류

12. Activate - 회원 가입 신청 활성화/비활성화

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/member/activate | member manager |

    - Request
        - activate: (boolean) 활성화 여부 (활성화: true, 비활성화: false)

    - Request Body example
    ```json
    {
        "activate": true
    }
    ```

    - Response
        - data.active: (boolean) 활성화 여부
        - error: (string) 에러 메시지 (활성화/비활성화 성공 시 empty)

    - Response Body example
    ```json
    {
        "data": {
            "active": true
        },
        "error": "Already active"
    }
    ```

    - Status Code
        - 200 OK: 활성화/비활성화 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 이미 활성화/비활성화 돼있는 경우, 시스템 오류

13. Graduates - 졸업자 목록 조회 (추후 졸업자 일괄 메일 발송 시 사용)

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/graduates | member manager |

    - Response
        - data.graduates: (Array&lt;JSON&gt;) 졸업자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "graduates": [
                    {
                        "id": "20190000",
                        "password": "asdf1234",
                        "name": "김희동",
                        "department": "예술대학 도자기학과",
                        "phone": "010-1234-5678",
                        "email": "heedong@gmail.com",
                        "grade": 4,
                        "attendance": 2,
                        "approved": true,
                        "on_delete": false,
                        "created_at": "1629000020",
                        "updated_at": "1629002720",
                        "role": {
                            "member_management": false,
                            "activity_management": false,
                            "fee_management": false
                        }
                    },
                    {
                        "id": "20200299",
                        "password": "asfjsk23242",
                        "name": "이기철",
                        "department": "공과대학 전자공학부",
                        "phone": "010-9876-5432",
                        "email": "lee@naver.com",
                        "grade": 4,
                        "attendance": 2,
                        "approved": true,
                        "on_delete": false,
                        "created_at": "1619060720",
                        "updated_at": "1619061720",
                        "role": {
                            "member_management": true,
                            "activity_management": true,
                            "fee_management": true
                        }
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 500 Internal Server Error: 시스템 오류