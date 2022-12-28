# WIP

## 개요
- 여러 패턴을 학습해본다.
- 프로젝트 레이아웃은 [golang-standards/project-layout](https://github.com/golang-standards/project-layout) 을 참고함
- sqs 를 사용한 이벤트 처리
- redis 를 사용한 응답 캐싱
- gorm 과 sqlc을 써봄. (어떤 구현체를 사용할지는 config.json에서 주입)
- 서비스레이어에서 트랜잭션 처리가 가능하도록 함 (별로 깔끔하게 작성되지는 않음)
- gracefulShutdown

## config
```json
{
  "db": {
    "db_port": "3306",
    "db_url": "your-db-url",
    "db_password": "your-db-password",
    "db_username": "your-db-username",
    "db_schema": "your-db-schema",
    "db_client": "sqlc | gorm"
  },
  "redis": {
    "redis_port": "6379",
    "redis_url": "your-redis-url",
    "redis_db": "db-for-caching-response"
  },
  "s3": {
    "s3_url": "",
    "s3_bucket_name": ""
  },
  "sqs_queue": {
    "queue_url": "your-sqs-url" 
  }
}
```
- s3를 이용한 기능은 아직 구현하지않음
- aws credential 파일 저장되어있어야 함 `(~/.aws/credentials)`

## 기타
- README 도 업뎃할 것이고, 무튼 자세한 구조를 설명할 예정..
- 재미 및 공부삼아 말도안되게 구현한 것들이 많음 (일부로 오래걸리는 기능인척 하려고 sleep 을 한다거나..)
- 서비스레이어에서 트랜잭션을 컨트롤하기 위해 컨텍스트에 트랜잭션 여부를 집어넣음
  - 예를들어 Transactional 하게 동작해야하는 서비스 메서드라면 context.WithValue 에 Transactional{} 을 true 로 넘기게 되면, repository layer 에서 트랜잭션을 사..
  - 이거 쫌 이상한것 같은데 여러 sql 구현체 (gorm, sqlc, ent 등등...) 을 동시에 사용하면서 서비스레이어에서 트랜잭션을 통제하려면 이것 말고 또 어떤 방법이 있는지 모르겠다..
  - 일단 context 자체를 cancel 시그널 보내는 용도로만 쓰는게 좋다고 한다 [출처](https://dave.cheney.net/2017/01/26/context-is-for-cancelation)
  - 그리고 얘때매 서비스 테스트용 mock repository 사용하는게 매우 귀찮아진다.
  - 이건 어떻게 하는게 좋을지.. 좀 더 생각을 해봐야한다.
- 근데 `service/task/task.go` `Create`메서드에서 트랜잭션을 사용하는 것은 사실 의미가 없음. 제데로 작동하는 예제 메서드를 만들어봐야 할듯
  - `create table ...` 구문은 implicit commit 으로 작동하기 때문에.... (저렇게 코딩해둔 뒤 알았음ㅋ)