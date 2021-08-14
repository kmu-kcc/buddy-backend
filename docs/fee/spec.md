# Buddy Back-end Fee API Specification

0. Server Domain:Port

    localhost:3000 (추후 변경 예정)

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
            "year": "2021",
            "semester": "2",
            "amount": "40000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회비 내역 생성 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

2. Amount - 회비 납부 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/create | manager |

    - Request
        - member_id: (string) 학번
        - year: (number) 연도
        - semester: (number) 학기

    - Request Body example
        ```json
        {
            "member_id": "20190000",
            "year": 2021,
            "semester": 2,
        }
        ```

    - Response
        - sum: (number) 총 납부한 회비 금액
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "data": {
                "sum": 20000,
            },
            "error": ""
        }
        ```

3. Dones - 회비 납부자 명단 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/fee/dones | manager|
    
    - Request
        - year: (number) 납부자를 조회할 연도
        - semester: (number) 납부자를 조회할 학기

    - Request Body example
        ```json
        {
            "year": "2021",
            "semester": "1"
        }


    - Response
        - dones: (Array&lt;JSON&gt;) 회비 납부자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "data": {
                "dones": [
                    {
                        "id": "20172229",
                        "name": "홍길동",
                        "department": "나노전자물리학과",
                        "grade": "1",
                        "phone": "010-2021-0001",
                        "email": "testmail1",
                        "attendance": 0
                    },
                    {
                        "id": "20171718",
                        "name": "심청이",
                        "department": "정보보안암호수학과",
                        "grade": "1",
                        "phone": "010-2021-0001",
                        "email": "testmail2",
                        "attendance": 0
                    }
                ]
            },
            "error": ""
        }
        ```

4. Yets - 회비 미납자 명단 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/fee/yets | manager |
    
    - Request
        - year:(number) 미납자를 조회할 연도
        - semester:(number) 미납자를 조회할 학기
        
    - Request Body example
        ```json
        {
            "year": "2021",
            "semester": "1"
        }

    - Response
        - dones: (Array&lt;JSON&gt;) 미납자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "data": {
                "yets": [
                    {
                        "id": "20172229",
                        "name": "홍길동",
                        "department": "나노전자물리학과",
                        "grade": "1",
                        "phone": "010-2021-0001",
                        "email": "testmail1",
                        "attendance": 0
                    },
                    {
                        "id": "20171718",
                        "name": "심청이",
                        "department": "정보보안암호수학과",
                        "grade": "1",
                        "phone": "010-2021-0001",
                        "email": "testmail2",
                        "attendance": 0
                    }
                ]
            },
            "error": ""
        }
        ```

5. All - 회비 내역 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/fee/all | member |
    
    - Request
        - startdate: (number) 회비 내역을 조회할 시작 기준 날짜
        - enddate: (number) 회비 내역을 조회할 끝 기준 날짜
        
    - Request Body example
        ```json
        {
            "startdate": 1627794157,
            "enddate": 1627794157
        }

    - Response
        - logs: (Array&lt;JSON&gt;) 회비 내역
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "data": {
                "logs": [
                      {
                        "id": "610521781cc9c3cc51f7c06c",
                        "member_id": "20172229",
                        "amount": "15000",
                        "type": "approved",
                        "updated_at": "0"
                      },
                      {
                        "id": "610521791cc9c3cc51f7c06d",
                        "member_id": "20190002",
                        "amount": "10000",
                        "type": "approved",
                        "updated_at": "0"
                      }
                ]
            },
            "error": ""
        }
        ```

6. Approve - 회비 납부 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/approve | manager |

    - Request
        - ids : (&lt;Array&gt;string) 납부 요청 목록
    
    - Request Body example
        ```json
        {
            "ids": [
                "610bf6b09a38451598148a25",
                "610bf6b09a38451598148a55"
            ]
        }
        ```
    
    - Response
        - data : (string) 반환 정보 (반환이 없을 경우 "")
        - error: (string) 에러 메시지 (정상 처리 시 "")
    
    - Response Body example
        ```json
        {
            "data" : [
                {
                    "Response": "Contents"
                }
            ],
            "error": "mongo: no such documents"
        }
        ```

7. Deposit - 입금/지출 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/deposit | manager |

    - Request
        - year : (string) 연도
        - semester : (string) 학기
        - amount : (string) 금액

    - Request Body example
        ```json
        {
            "year" : "2021",
            "semester" : "4",
            "amount" : "10000"
        }
        ```
    
    - Response
        - data : (string) 반환 정보 (반환이 없을 경우 "")
        - error: (string) 에러 메시지 (정상 처리 시 "")
    
    - Response Body example
        ```json
        {
            "data" : [
                {
                    "Response": "Contents"
                }
            ],
            "error": "mongo: no such documents"
        }
        ```