# Buddy Backend Fee API Specification

0. Server Domain:Port

    localhost:3000 (추후 변경 예정)

<br>

1. Create - 회비 내역 생성

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/create | manager |

    - Request
        - year: (number) 연도
        - semester: (number) 학기
        - amount: (number) 금액

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
        
2. Submit - 회비 납부 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/submit | member |

    - Request
        - member_id: (string) 학번
        - year: (number) 연도
        - semester: (number) 학기
        - amount: (number) 금액

    - Request Body example
        ```json
        {   
            "member_id": "20190000",
            "year": "2021",
            "semester": "2",
            "amount": "20000"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (회비 납부 신청 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```
        
3. Amount - 회비 납부 조회

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

6. Dones - 회비 납부자 명단 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/fee/dones | manager|
    
    - Request
        - year: (number) 납부자를 조회할 연도
        - semester: (number) 납부자를 조회할 학기

    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 1
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

7. Yets - 회비 미납자 명단 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/fee/yets | manager |
    
    - Request
        - year:(number) 미납자를 조회할 연도
        - semester:(number) 미납자를 조회할 학기
        
    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 1
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
10. All - 회비 내역 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/fee/all | member |
    
    - Request
        - year: (number) 회비 내역을 조회할 연도
        - semester: (number) 회비 내역을 조회할 학기
        
    - Request Body example
        ```json
        {
            "year": 2021,
            "semester": 1
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

11. Approve - 회비 납부 요청 승인

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
        - error: (string) 에러 메시지 (정상 처리 시 "")
    
    - Response Body example
        ```json
        {
            "error": "mongo: no such documents"
        }
        ```    

12. Reject - 회비 납부 처리 (거부)

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/fee/reject | manager |

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
        - error: (string) 에러 메시지 (정상 처리 시 "")
    
    - Response Body example
        ```json
        {
            "error": "mongo: no such documents"
        }
        ```    

13. Deposit - 입금/지출 처리

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
        - error: (string) 에러 메시지 (정상 처리 시 "")
    
    - Response Body example
        ```json
        {
            "error": "mongo: no such documents"
        }
        ```
