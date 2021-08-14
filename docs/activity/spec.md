# Buddy Back-end Activity API Specification

0. Server Domain:Port

    localhost:3000 (추후 변경 예정)

<br>

1. Create - 활동 생성

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/create | manager |

    - Request
        - start: (number) 시작일
        - end: (number) 종료일
        - place: (string) 장소
        - type: (string) 종류
        - description: (string) 설명
        - participants: (Array&lt;string&gt;) 참여자 목록
        - private: (bool) public/private

    - Request Body example
        ```json
        {
            "start": 1,
            "end": 1,
            "place": "cafe",
            "type": "study",
            "description": "good",
            "participants" : [],
            "private": true
        }
        ```

    - Response
        - error: (string) 에러 메시지 (활동 생성 성공 시 "")

    - Response Body example
        ```json
        {
            "error": ""
        }
        ```


2. Search - 활동 검색

    | method | route | priviledge |
    | :---: | :---: | :---: |
    | POST | /api/v1/activity/search | manager |

    - Request
        - query: (string) 검색어

    - Request Body example
        ```json
        {
            "query": "st"
        }
        ```

    - Response
        - activities: (Array&lt;string&gt;) 활동 검색 결과
        - error: (string) 에러 메시지 (활동 검색 성공 시 "")

    - Response Body example
        ```json
        {
            "data": {
                "activities": [
                    {
                        "start": "1628249722",
                        "end": "1628249722",
                        "place": "cafe",
                        "type": "study",
                        "description": "Study start!",
                        "participants": [
                            "20192019",
                            "20182018"
                        ]
                    }
                ]
            },
            "error": ""
        }
        ```

3. Update - 활동 정보 갱신

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
                "participants": [
                        "20192019",
                        "20182018"
                    ]
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

4. Delete - 활동 삭제 처리

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