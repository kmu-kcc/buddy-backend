# Buddy Backend Activity API Specification

0. Server Domain:Port

    localhost:3000 (추후 변경 예정)

<br>

1. Search - 활동 검색

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/search | member |

    - Request
        - filter: (JSON) 검색하고자 하는 활동 정보 (시작일, 종료일, 장소, 종류, 설명, 참여자 중 0개 이상 택)

    - Request Body example
        ```json
        {
            "filter": {
                "start": "1628249722",
                "end": "1628249722",
                "place": "cafe",
                "type": "study",
                "description": "Study",
                "participants": "홍길동"
            }
        }
        ```

    - Response
        - activities: (Array&lt;string&gt;) 활동 검색 결과
        - error: (string) 에러 메시지 (활동 검색 성공 시 "")

    - Response Body example
        ```json
        {
            "activities": [
                {
                    "start": "1628249722",
                    "end": "1628249722",
                    "place": "cafe",
                    "type": "study",
                    "description": "Study start!",
                    "participants": "홍길동, 김철수"
                }
            ],
            "error": ""
        }
        ```
        
2. Update - 활동 정보 갱신

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/update | manager |

    - Request
        - update: (JSON) 갱신하고자 하는 활동 정보 (시작일, 종료일, 장소, 종류, 설명, 참여자 중 0개 이상 택)

    - Request Body example
        ```json
        {
            "update": {
                "_id": "610d458b79e122ea1d150cd6",
                "start": "1628249722",
                "end": "1628249722",
                "place": "cafe",
                "type": "study",
                "description": "Study End!",
                "participants": "홍길동"
            }
        }
        ```

    - Response
        - error: (string) 에러 메시지 (활동 정보 갱신 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```
        
3. Delete - 활동 삭제 처리

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/delete | manager |

    - Request
        - _id: (string) 삭제할 활동 id

    - Request Body example
        ```json
        {
            "_id": "610d458b79e122ea1d150cd6"
        }
        ```

    - Response
        - error: (string) 에러 메시지 (활동 삭제 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```
        
4. Participants - 활동 참여자 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/participants | manager |

    - Request
        - _id: (string) 참여자 목록을 볼 활동 id 

    - Request Body example
        ```json
        {
            "_id": "610d458b79e122ea1d150cd6"
        }
        ```

    - Response
        - members: (Array&lt;string&gt;) 활동 참여자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "members": [
                {
                    "id": "20180000",
                    "name": "홍길동",
                    "department": "스포츠레저학과",
                    "phone": "010-1234-5678",
                    "email": "gildong@gmail.com",
                    "grade": 4,
                    "attendance": 0
                },
                {
                    "id": "20190000",
                    "name": "김철수",
                    "department": "국어국문학과",
                    "phone": "010-5678-1234",
                    "email": "chulsoo@gmail.com",
                    "grade": 3,
                    "attendance": 0
                }
            ],
            "error": ""
        }
        ```

5. ApplyP - 활동 지원

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/activity/applyp | member |

    - Request
        - _id : (string) 활동 고유 ID
        - member_id : (string) 활동 신청 회원 ID

    - Request Body example
        ```json
        {
            "_id": "610d4cfa567af2cc318e7f97",
            "member_id": "20172227"
        }

    - Response
        - error: (string) 에러 메시지 (활동 신청 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

6. Papplies - 활동 지원자 목록 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | GET | /api/v1/activity/papplies | manager |

    - Request
        - _id : (string) 활동 고유 ID

    - Request Body example
        ```json
        {
            "_id": "610d4cfa567af2cc318e7f97"
        }

    - Response
        - papplies: (Array&lt;JSON&gt;)활동 지원자 목록
        - error: (string) 에러 메시지 (쿼리 성공 시 "")

    - Response Body example
        ```json
        {
            "data": {
                "papplies": [
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
                        "name": "홍길동2",
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

7. ApproveP - 활동 참여 승인

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/activity/approvep | manager |
    
    - Request
        - _id: (string )활동 고유 ID
        - member_ids: (Array&lt;string&gt;)참여 승인 회원 ID
        
    - Request Body example
        ```json
        {
            "_id": "610d4cfa567af2cc318e7f97",
            "member_ids": [
                "20172227",
                "20171718"
            ]
        }

    - Response
        - error: (string) 에러 메시지 (참여 승인 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

8. RejectP - 활동 참여 거절

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/activity/rejectp | manager |
    
    - Request
        - _id : (string) 활동 고유 ID
        - member_ids : (Array&lt;string&gt;) 참여 신청을 거절할 회원들의 ID List
        
    - Request Body example
        ```json
        {
            "_id": "610d4cfa567af2cc318e7f97",
            "member_ids": [
                "20172227",
                "20171718"
            ]
        }

    - Response
        - error: (string) 에러 메시지 (활동 참여 거절 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

9. CancelP - 활동 참여 신청 취소

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | PUT | /api/v1/activity/cancelp | member |
    
    - Request
        - _id : (string) 활동 고유 ID
        - member_id : (string) 활동 참여를 취소 신청하는 회원의 ID

    - Request Body example
        ```json
        {
            "_id": "610d4cfa567af2cc318e7f97",
            "member_id": "20172227"
        }

    - Response
        - error: (string) 에러 메시지 (참여 신청 취소 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```

10. ApplyC - 활동 참여 취소 신청

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/applyc | member |

    - Request
        - activity_id: (string) 활동 ID
        - member_id: (string) 신청자 ID
    
    - Request Body example
        ```json
        {
            "activity_id": "610bf6b09a38451598148a25",
            "member_id": "610bf6b09a38451598148a55",
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

11. CancelC - 활동 참여 취소 신청 취소

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/cancelc | member |

    - Request
        - activity_id: (string) 활동 ID
        - member_id: (string) 신청자 ID
    
    - Request Body example
        ```json
        {
            "activity_id": "610bf6b09a38451598148a25",
            "member_id": "610bf6b09a38451598148a55",
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

12. Capplies - 취소 신청자 리스트 조회

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/capplies | manager |

    - Request
        -
    
    - Request Body example
    
    - Response
        - participants: (Array&lt;string&gt;) 참가자 리스트
        - error: (string) 에러 메시지 (정상 처리 시 "")
    
    - Response Body example
        ```json
        {
            "participants" : [
                "20000808"
                "19990803"
            ],
            "error": "mongo: no such documents"
        }
        ```    

13. ApproveC - 참여 취소 신청 승인

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/approvec | manager |

    - Request
        - activity_id: (string) 활동 ID
        - member_id : (string) 멤버 ID
    
    - Request Body example
        ```json
        {
            "activity_id": "690bf6b09a38451598148a25",
            "member_id": "610bf6b09a38451598148a55"
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

14. RejectC - 참여 취소 신청 거부

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/rejectc | manager |

    - Request
        - activity_id: (string) 활동 ID
        - member_id : (string) 멤버 ID
    
    - Request Body example
        ```json
        {
            "activity_id": "610bf6b09a38451598148a25",
            "member_id": "610bf6b09a38451598148a55"
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
