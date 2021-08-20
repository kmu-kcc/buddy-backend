# Buddy Back-end Fee API Specification

0. Server Domain:Port

    http://localhost:3000 (추후 변경 예정)

<br>

1. Create - 회비 내역 초기화

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/create | manager |

    - Request
        - year: (number) 연도
        - semester: (number) 학기
        - amount: (number) 해당 학기에 1인당 납부해야할 금액

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 2,
            "amount": 15000
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회비 내역 초기화 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "duplicated fee"
        }
        ```

    - Status code
        - 200 OK: 회비 내역 초기화 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 회비 내역 중복, 시스템 오류 등

2. Amount - 회비 납부액 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/amount | member |

    - Request
        - id: (string) 학번
        - year: (number) 연도
        - semester: (number) 학기

    - Request Body example
        ```json
        {
            "id": "20190000",
            "year": 2021,
            "semester": 2,
        }
        ```

    - Response
        - data.amount: (number) 해당 학기에 납부한 총 회비 금액
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "amount": 20000,
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

3. Payers - 회비 납부자 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/payers | manager |
    
    - Request
        - year: (number) 조회할 연도
        - semester: (number) 조회할 학기

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 1
        }

    - Response
        - data.payers: (Array&lt;JSON&gt;) 회비 납부자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "payers": [
                    {
                        "id": "20172229",
                        "password": "asdf2322",
                        "name": "홍길동",
                        "department": "공과대학 나노전자물리학과",
                        "phone": "010-2021-0001",
                        "email": "gildong@kookmin.ac.kr",
                        "grade": 1,
                        "attendance": 0,
                        "approved": true,
                        "on_delete": false,
                        "created_at": "1629060720",
                        "updated_at": "1629060720"
                    },
                    {
                        "id": "20171718",
                        "password": "8809dfsfdsf",
                        "name": "심청이",
                        "department": "공과대학 정보보안암호수학과",
                        "phone": "010-2021-0001",
                        "email": "simch@naver.com",
                        "grade": 1,
                        "attendance": 0,
                        "approved": true,
                        "on_delete": false,
                        "created_at": "1629060720",
                        "updated_at": "1629060720"
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

4. Deptors - 회비 미납자 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/deptors | manager |
    
    - Request
        - year:(number) 조회할 연도
        - semester:(number) 조회할 학기
        
    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 1
        }

    - Response
        - data.deptors: (Array&lt;JSON&gt;) 미납자 정보 및 미납액 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "deptors": [
                    {
                        "id": "20172229",
                        "password": "asdf2322",
                        "name": "홍길동",
                        "department": "공과대학 나노전자물리학과",
                        "phone": "010-2021-0001",
                        "email": "gildong@kookmin.ac.kr",
                        "grade": 1,
                        "attendance": 0,
                        "approved": true,
                        "on_delete": false,
                        "created_at": "1629060720",
                        "updated_at": "1629060720",
                        "dept": 3000
                    },
                    {
                        "id": "20171718",
                        "password": "8809dfsfdsf",
                        "name": "심청이",
                        "department": "공과대학 정보보안암호수학과",
                        "phone": "010-2021-0001",
                        "email": "simch@naver.com",
                        "grade": 1,
                        "attendance": 0,
                        "approved": true,
                        "on_delete": false,
                        "created_at": "1629060720",
                        "updated_at": "1629060720",
                        "dept": 15000
                    }
                ]
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

    - Status Code
        - 200 OK: 쿼리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

5. Search - 회비 내역 검색

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/search | member |
    
    - Request
        - year: (number) 조회할 연도
        - semester: (number) 조회할 학기

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 1
        }

    - Response
        - data.carry_over: (number) 이월 금액
        - data.logs: (Array&lt;JSON&gt;) 회비 내역 (`type` - 회비 납부: 0, 입/출금: 1)
        - data.total: (number) 계
        - error: (string) 에러 메시지 (쿼리 성공 시 empty)

    - Response Body example
        ```json
        {
            "data": {
                "carry_over": 150000,
                "logs": [
                    {
                        "amount": 15000,
                        "type": 0,
                        "description": "회비 납부",
                        "created_at": "1619060720"
                    },
                    {
                        "amount": -10000,
                        "type": 1,
                        "description": "비품 구입",
                        "created_at": "1629000020"
                    }
                ],
                "total": 155000
            },
            "error": "argument to Unmarshal* must be a pointer to a type, but got ..."
        }
        ```

6. Pay - 회비 납부 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/pay | manager |

    - Request
        - year: (number) 납부 처리할 연도
        - semester: (number) 납부 처리할 학기
        - payments : (Array&lt;JSON&gt;) 납부 처리 목록

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 2,
            "payments": [
                {
                    "id": "20189879",
                    "amount": 15000
                },
                {
                    "id": "20209013",
                    "amount": 15000
                },
                {
                    "id": "20170907",
                    "amount": 10000
                }
            ]
        }
        ```
    
    - Response
        - error: (string) 에러 메시지 (납부 처리 성공 시 empty)
    
    - Response Body example
        ```json
        {
            "error": "mongo: no such documents"
        }
        ```

    - Status Code
        - 200 OK: 회비 납부 처리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

7. Deposit - 입금/지출 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/deposit | manager |

    - Request
        - year: (number) 연도
        - semester: (number) 학기
        - amount: (number) 금액 (입금일 경우 양수, 지출일 경우 음수)
        - description: (string) 비고

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 2,
            "amount": 100000,
            "description": "운영비 지원"
        }
        ```
    
    - Response
        - error: (string) 에러 메시지 (입금 처리 성공 시 empty)
    
    - Response Body example
        ```json
        {
            "error": "mongo: no such documents"
        }
        ```

    - Status Code
        - 200 OK: 입금 처리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 시스템 오류

8. Exempt - 면제 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/exempt | manager |

    - Request
        - year: (number) 연도
        - semester: (number) 학기
        - id: (string) 회비 면제 대상의 학번

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 2,
            "id": "20210001"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (면제 처리 성공 시 empty)

    - Response Body example
        ```json
        {
            "error": "already exempted"
        }
        ```

    - Status Code
        - 200 OK: 면제 처리 성공
        - 400 Bad Request: 요청 포맷/타입 오류
        - 500 Internal Server Error: 이미 면제 처리된 경우, 시스템 오류