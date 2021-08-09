# Buddy Backend Member API Specification

0. Server Domain:Port

    localhost:3000 (추후 변경 예정)

<br>

1. SignIn - 회원 로그인

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
        - error: (string) 에러 메시지 (로그인 성공 시 "")

    - Response Body example
        ```json
        {
            "error": "password mismatch"
        }
        ```

2. SignUp - 회원 가입 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/signup | - |

    - Request
        - id: (string) 학번
        - name: (string) 이름
        - department: (string) 소속 학부
        - phone: (string) 전화번호
        - email: (string) 메일 주소
        - grade: (number) 학년 (1 이상의 정수)
        - attendance: (number) 재학 여부 (재학: 0, 휴학: 1: 졸업: 2)

    - Request Body example
        ```json
        {
            "id": "20210000",
            "name": "홍길동",
            "department": "시각디자인학과",
            "phone": "010-1234-5678",
            "email": "gildong@gmail.com",
            "grade": 1,
            "attendance": 0
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회원 가입 신청 성공 시 "")

    - Response Body example
        ```json
        {
            "error": "under review"
        }
        ```

3. SignUps - 회원 가입 신청 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/signups | manager |

    - Response
        - signups: (Array&lt;JSON&gt;) 회원 가입 신청자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "signups": [
                {
                    "id": "20190000",
                    "name": "김희동",
                    "department": "스포츠레저학과",
                    "phone": "010-1234-5678",
                    "email": "heedong@gmail.com",
                    "grade": 3,
                    "attendance": 0
                },
                {
                    "id": "20200299",
                    "name": "이기철",
                    "department": "물리학과",
                    "phone": "010-9876-5432",
                    "email": "lee@naver.com",
                    "grade": 2,
                    "attendance": 1
                }
            ],
            "error": ""
        }
        ```

4. Approve - 회원 가입 승인

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/approve | manager |

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
        - error: (string) 에러 메시지 (회원 가입 승인 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

5. Exit - 회원 탈퇴 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/exit | member |

    - Request
        - id: (string) 탈퇴 신청하는 회원의 학번

    - Request Body example
        ```json
        {
            "id": "20210000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (탈퇴 신청 성공 시 "")

    - Response Body example
        ```json
        {
            "error": "already on delete"
        }
        ```

6. Exits - 회원 탈퇴 신청 목록

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/exits | manager |

    - Response
        - exits: (Array&lt;JSON&gt;) 회원 탈퇴 신청자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "exits": [
                {
                    "id": "20190000",
                    "name": "김희동",
                    "department": "스포츠레저학과",
                    "phone": "010-1234-5678",
                    "email": "heedong@gmail.com",
                    "grade": 3,
                    "attendance": 0
                },
                {
                    "id": "20200299",
                    "name": "이기철",
                    "department": "물리학과",
                    "phone": "010-9876-5432",
                    "email": "lee@naver.com",
                    "grade": 2,
                    "attendance": 1
                }
            ],
            "error": ""
        }
        ```

7. Delete - 회원 가입 승인 거부 및 회원 탈퇴 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/delete | manager |

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
        - error: (string) 에러 메시지 (회원 가입 승인 거부/탈퇴 처리 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

8. CancelExit - 회원 탈퇴 신청 취소

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/cancelexit | member |

    - Request
        - id: (string) 탈퇴 신청을 취소하는 회원의 학번

    - Request Body example
        ```json
        {
            "id": "20210000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (탈퇴 신청 취소 성공 시 "")

    - Response Body example
        ```json
        {
            "error": "not on delete"
        }
        ```

9. Search - 회원 검색

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/search | member |

    - Request
        - filter: (JSON) 검색하고자 하는 회원 정보 (학번, 이름, 소속 학부, 학년 중 0개 이상 택)

    - Request Body example
        ```json
        {
            "filter": {
                "id": "20210000",
                "name": "홍길동",
                "department": "소프트웨어학부",
                "grade": 1
            }
        }
        ```

    - Response
        - members: (Array&lt;JSON&gt;) 회원 검색 결과
        - error: (string) 에러 메시지 (회원 검색 성공 시 "")

    - Response Body example
        ```json
        {
            "members": [
                {
                    "id": "20210000",
                    "name": "홍길동",
                    "department": "소프트웨어학부",
                    "email": "gildong@gmail.com",
                    "grade": 1
                }
            ],
            "error": ""
        }
        ```

10. Update - 회원 정보 갱신

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/update | member |

    - Request
        - update: (JSON) 갱신하고자 하는 회원 정보 (비밀번호, 이름, 소속 학부, 전화번호, 이메일, 학년, 재학 여부 중 0개 이상 택)

    - Request Body example
        ```json
        {
            "update": {
                "password": "asdf1234",
                "name": "홍길동",
                "department": "소프트웨어학부",
                "phone": "010-1234-5678",
                "email": "gildong@yahoo.com",
                "grade": 1,
                "attendance": 1
            }
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회원 정보 갱신 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

11. ApplyGraduate - 졸업 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/applygraduate | member |

    - Request
        - id: (string) 졸업 신청는 회원의 학번

    - Request Body example
        ```json
        {
            "id": "20210000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (졸업 신청 성공 시 "")

    - Response Body example
        ```json
        {
            "error": "already on graduate"
        }
        ```

12. GraduateApplies - 졸업 신청자 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/graduateapplies | manager |

    - Response
        - applies: (Array&lt;JSON&gt;) 졸업 신청자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "applies": [
                {
                    "id": "20190000",
                    "name": "김희동",
                    "department": "스포츠레저학과",
                    "phone": "010-1234-5678",
                    "email": "heedong@gmail.com",
                    "grade": 4,
                    "attendance": 0
                },
                {
                    "id": "20200299",
                    "name": "이기철",
                    "department": "물리학과",
                    "phone": "010-9876-5432",
                    "email": "lee@naver.com",
                    "grade": 4,
                    "attendance": 0
                }
            ],
            "error": ""
        }
        ```

13. CancelGraduate - 졸업 신청 취소

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/cancelgraduate | member |

    - Request
        - id: (string) 졸업 신청을 취소하는 회원의 학번

    - Request Body example
        ```json
        {
            "id": "20210000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (졸업 신청 취소 성공 시 "")

    - Response Body example
        ```json
        {
            "error": "not on graduate"
        }
        ```

14. ApproveGraduate - 졸업 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/member/approvegraduate | manager |

    - Request
        - ids: (Array&lt;string&gt;) 졸업 처리하는 회원들의 학번 List

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
        - error: (string) 에러 메시지 (졸업 처리 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

15. Graduates - 졸업자 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/member/graduateapplies | manager |

    - Response
        - graduates: (Array&lt;JSON&gt;) 졸업자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "graduates": [
                {
                    "id": "20190000",
                    "name": "김희동",
                    "department": "스포츠레저학과",
                    "phone": "010-1234-5678",
                    "email": "heedong@gmail.com",
                    "grade": 4,
                    "attendance": 2
                },
                {
                    "id": "20200299",
                    "name": "이기철",
                    "department": "물리학과",
                    "phone": "010-9876-5432",
                    "email": "lee@naver.com",
                    "grade": 4,
                    "attendance": 2
                }
            ],
            "error": ""
        }
        ```
