# Event-Driven Architecture

## Introduction

회사 서비스가 성장하고 요구사항이 다양해 짐에 따라 기존의 모놀리식 아키텍처에 한계를 느꼈습니다. 그리고 올해(2025년), 위 문제를 해결하기 위해 모놀리식 아키텍처인 서비스를 MSA로 전환하는 중입니다.

이 과정에서 각 서비스에서 발생하는 이벤트를 다른 서비스에 적용하기 위해 여러 삽질?(REST API 요청, 백그라운드 작업, Redis PubSub...)을 시도했습니다. 구현하면서도 이 방식이 멋지지 않고 전혀 MSA스럽지 않음을 인지하고 있었지만, 서비스의 규모를 생각했을 때 오버 엔지니어링이라는 생각에 EDA 적용을 미뤘습니다(물론 흐린 눈으로 바라 본 것도 없지 않아 있습니다..).

하지만 이렇게 계속 미룰 수는 없기에 공부해 보려 합니다. 이 문서는 다음과 같은 의문 및 궁금증을 해결하기 위해 EDA를 공부하며 정리한 문서입니다:

- 진짜 EDA 적용이 필요한 것이 맞을까?
- 대체 어느 정도 규모가 되어야 EDA 적용이 오버 엔지니어링이 아닐까?
- 회사에 EDA 적용을 어떻게 건의해야 할까?

## Event-Driven Architecture (EDA) 란?

- 이벤트 기반 아키텍처(EDA)는 시스템 내 발생하는 이벤트를 기반으로 컴포넌트(서비스)들이 통신하는 구조를 의미함
    - 컴포넌트(서비스)들이 이벤트를 생성하고 이에 응답함으로써 서로 통신하는 소프트웨어 아키텍처임
- EDA는 이벤트를 사용하여 분리된 서비스 간의 통신 및 상호 작용을 트리거하며, [마이크로서비스](../micro_service_acrhitecture/MICRO_SERVICE_ARCHITECTURE.md)로 구축된 시스템에서 흔히 사용됨
    - 특정 이벤트가 발생하면, 이를 구독하고 있는 다른 컴포넌트들이 이에 반응하여 동작함

### 이벤트란?

- 이벤트는 상태 변화 또는 업데이트를 의미함
    - e.g. 장바구니에 들어간 품목, 스토리지 시스템에 업데이트된 파일, 배송 준비된 주문
- 이벤트는 상태 정보(e.g. 주문 내의 품목 이름, 가격, 수량 등)를 전달할 수도 있고, 관련 정보를 조회하는 데 필요한 식별자(e.g. `주문 번호 1234가 배송되었습니다.`)만을 포함하기만 할 수도 있음

#### 이벤트의 주요 특징

- Representation (표현): 이벤트는 특정 정보를 전달하는 메시지나 신호의 형태로 표현됨
- Triggering: 사용자 작업이나 데이터 변경과 같은 다양한 소스에 의해 이벤트가 발생함
- Asynchronicity (비동기성): EDA는 주로 비동기 통신을 사용하여 각 컴포넌트가 서로 독립적이고, 병렬적으로 작동할 수 있도록 합니다.
- Publish-Subscribe Model: 이벤트 관리를 위해 Publish-Subscribe Model이 사용됨
- Event Types: 이벤트는 그 목적에 따라 그룹화 됨(e.g. `UserLoggedIn`, `OrderPlaced`)
- Payload: 이벤트에는 종종 맥락을 제공하는 추가 정보인 'payload'가 포함됨
- Event Handling: 각 컴포넌트에는 이벤트에 어떻게 반응할지 결정하는 전용 핸들러가 있음
- Real-Time Processing: 이벤트를 통해 변화에 즉각적으로 반응할 수 있으므로, 빠른 응답성이 요구되는 시나리오에서 EDA는 매우 이상적임

### EDA의 특징

- 비동기성
- 느슨한 결합

### Event-Driven Architecture의 구성 요소

#### Event Producer (이벤트 생산자)

- 이벤트를 생성하여 이벤트 브로커로 전송하는 구성요소임
- 이벤트 생산자는 특정 조건이 충족되거나 작업이 발생할 때 이벤트를 방출(emit)하는 책임을 짐
- 이벤트 소비자를 특정하지 않음

#### Event Broker (이벤트 브로커)

- 이벤트 브로커는 중앙 허브 역할을 하며, 이벤트 분배, 필터링 및 라우팅을 처리함으로써 다양한 컴포넌트 간의 통신을 원할하게 함
    - 이벤트 생산자와 이벤트 소비자를 연결하는 핵심적인 구성 요소
- 다음과 같은 요구사항을 충족해야 함
    - 일종의 큐(Queue) 역할
    - 일정 기간 이벤트를 저장할 수 있는 기능 필요
    - 다수의 연결 (생산자, 소비자) 을 처리할 수 있는 확장성 필요
    - 단일 장애 포인트가 될 수 있기 때문에 오류 처리 기능(e.g. Dead Letter Queue) 필요

#### Event Consumer (이벤트 소비자)

- 특정 이벤트 유형에 관심을 보이고 이를 소비하는 구성 요소
- 이벤트 소비자는 이벤트 브로커에서 관련 이벤트를 대기(listen)하고 있다가, 이벤트가 발생하면 그에 따른 행동을 취함
- 특정 이벤트에 다수의 이벤트 소비자가 존재할 수 있음
    - 따라서 한 이벤트에 대해 여러 소비자가 작업을 시작하는 병렬 처리가 가능함

## Event-Driven Architecture이 필요성

그렇다면 왜 많은 서비스들이 Event-Driven Architecture를 선택할까?

Event-Driven Architecture의 적용으로 충족할 수 있는 요구조건에 대해 알아보고,
Event-Driven Architecture가 적용된 use case들을 알아 봄으로써 왜 많은 서비스들이 Event-Driven Architecture을 적용했는지 알아보자.

### Event-Driven Architecture이 필요한 경우

EDA는 MSA를 사용하는 애플리케이션이나, 컴포넌트가 분리된 애플리케이션에서 흔히 볼 수 있음.

특히, 다음과 같은 시나리오에 적합함:

- Integration of Heterogeneous Systems (이종 시스템의 통합)
    - 서로 통신해야 하는 다양한 시스템이나 서비스를 사용하는 경우
        - EDA는 서로 다른 기술 간 유연한 통신이 가능하게 함
        - 이벤트 브로커는 시스템 간 간접 연결 및 상호 운용성을 구축하여 시스템들이 기술 스택에 구애받지 않고 메시지와 데이터를 교환할 수 있도록 함
- Scalability Needs (확장성 요구)
    - 시스템이 성장하고 처리해야 할 이벤트 수가 증가할 것으로 예상되는 경우
        - EDA는 더 나은 확장성을 제공함
        - 전체 시스템을 중단하지 않고도 컴포넌트를 추가하거나 수정할 수 있음
- Fanout and Parallel processing (팬아웃 및 병렬 처리)
    - 이벤트에 반응하여 작동해야 하는 시스템 혹은 서비스가 많은 경우
        - EDA를 사용하면 각 시스템에 이벤트를 푸시하기 위해 사용자 지정 코드를 작성할 필요 없이 이벤트를 분산시킬 수 있음
        - 이벤트 브로커가 시스템에 이벤트를 푸시하면 각 시스템은 서로 다른 목적을 가지고 이벤트를 병렬로 처리할 수 있음
- Real-Time Applications (실시간 애플리케이션)
    - 애플리케이션이 사용자 동작이나 데이터 변경에 즉각적으로 반응해야 하는 경우
- Complex Event Processing(복합 이벤트 처리)
    - 여러 이벤트를 처리하고 이를 통해 인사이트를 도출해야 하는 애플리케이션인 경우
        - EDA를 통해 복잡한 시나리오 처리를 간소화할 수 있음
- Decoupled Components (분리된 구성 요소)

### Event-Driven Architecture의 Use Cases

- 마이크로서비스 간 통신
    - 전자상거래나 미디어 및 엔터테인먼트 플랫폼는 예측 불가능한 트래픽을 처리하기 위해 확장성이 요구되는 경우가 많음
        - 전자상거래 플랫폼에서 주문을 하면, 주문 이벤트가 이벤트 브로커로 전송됨
        - 하위 마이크로서비스들은 모두 이 주문 이벤트를 수신하여 처리할 수 있음
        - 예를 들어, 주문 접수, 재고 변동, 결제 처리 등의 작업을 수행할 수 있음
    - 각 마이크로서비스는 독립적으로 확장 및 장애 처리가 가능하므로, 단일 장애 지점 없이 주문량이 많은 시간대에도 프로세스를 확장할 수 있음
- 비즈니스 워크플로 자동화
    - 금융 서비스 거래와 같은 많은 비즈니서 워크플로는 동일한 단계를 반복해야 함
        - EDA를 사용하면 이러한 단계를 자동화할 수 있음
    - 또한, 거래 실행, 결제 승인 및 시장 변동과 같은 이벤트에 즉각적인 대응을 트리거할 수 있음
- SaaS 애플리케이션 통합
    - SaaS 환경의 가장 큰 과제는 사용자 활동 및 데이터에 대한 가시성 부족임
        - 사일로화된 데이터를 활용하기 위해 EDA는 SaaS 애플리케이션 이벤트를 수집하거나, SaaS 애플리케이션으로 이벤트를 전송할 수 있음
        - 예를 들어, 파트너 주문 데이터를 수집하고, 해당 주문을 사내 주문 처리 애플리케이션으로 직접 전송하는 미들웨어를 구축할 수 있음
- 인프라 자동화
    - 연산 집약적인 워크로드를 실행할 때, 고도의 병렬 처리를 위해 컴퓨팅 리소스를 확장하고, 작업 완료 후 축소하는 방식으로 대응하도록 설정할 수 있음
        - 예를 들어, 규제가 엄격한 산업 분야에서 EDA를 사용하면, 사고 발생 시 보안 태세 리소스를 대폭 확충하거나, 보안 정책에 따라 경고 이벤트가 발생할 때마다 시장 조치를 취할 수 있음
    - EDA를 활용하면 통신 산업에서도 네트워크 구성 요소 간의 이벤트 기반 통신, 네트워크 모니터링 및 통화 처리를 지원할 수 있음
        - 이를 통해 동적인 네트워크 환경 관리 및 부하 적응을 용이하게 함
- IoT
- Online Game

## Event-Driven Architecture의 장단점

### Event-Driven Architecture의 장점

- Flexibility and Agility (유연성과 반응성)
    - EDA는 컴포넌트를 분리함으로써 시스템의 변화하는 요구 사항에 신속하게 적응할 수 있도록 함
    - 전체 시스템에 영향을 주지 않고 새로운 기능을 추가하거나 변경할 수 있음
- Scalability (확장성)
    - EDA는 컴포넌트들이 독립적으로 작동할 수 있도록 함으로써 확장성을 지원함
    - 시스템은 컴포넌트나 리소스를 추가하여 증가하는 부하 또는 늘어나는 데이터 세트를 처리할 수 있음
- Real-time Processing (실시간 처리)
    - EDA는 실시간 처리가 필요한 시나리오에 이상적임
    - 이벤트가 발생하는 즉시 처리되므로 시스템은 시간에 민감한 작업을 효율적으로 관리할 수 있음
        - 병목 현상을 줄이고, 시스템 전체의 처리량을 높일 수 있음
- Loose Coupling (느슨한 결합)
    - 이벤트 기반 시스템의 컴포넌트들은 느슨하게 연결되어 있어 서로에게 크게 의존하지 않음
        - 이는 각 컴포넌트의 자립성을 높이고, 개별 구성 요소의 유지 관리를 간소화 함
- Enhanced Modulaarity
    - EDA는 복잡한 시스템을 관리 가능한 컴포넌트로 분해하는 모듈식 설계를 권장함
        - 이러한 모듈식 구조는 개발, 테스트 및 유지 관리를 간소화 함


### Event-Driven Architecture의 단점

- Increased Complexity (복잡성 증가)
    - 이벤트와 구성 요소가 추가될 수록 EDA 시스템은 복잡해질 수 있음
        - 이벤트 흐름을 관리하고 조율하는 것이 어려워질 수 있음
- Event Order and Consistency (이벤트 순서 및 일관성)
    - 이벤트를 올바른 순서로 유지하고, 시스템의 일관성을 보장하는 것이 까다로울 수 있음
    - 순서대로 발생하지 않는 이벤트를 처리하거나 작업이 그룹으로 완료되도록 하려면 추가적인 노력이 필요함
- Debugging and Tracing (디버깅 및 추적)
    - 분산 및 비동기 환경에서 문제를 찾아 수정하는 것은 기존 시스템보다 어려울 수 있음
        - 문제 해결에 더 많은 시간이 소요될 수 있음
- Event Latency (이벤트 지연 시간)
    - 이벤트는 개별적으로 처리되므로 이벤트 발생 시간과 응답 시점 사이에 지연이 발생할 수 있음
        - 이러한 지연은 빠른 반응이 필요한 상황에서 문제가 될 수 있음
- Eventual Consistency (최종 일관성)
- Observability (관측성)

## Event-Driven Architecture 주요 패턴

### Event Notification

- [EVENT_NOTIFICATION.md](./EVENT_NOTIFICATION.md)에서 정리 예정

### Event-Carried State Transfer

- [EVENT_CARRIED_STATE_TRANSFER.md](./EVENT_CARRIED_STATE_TRANSFER.md)에서 정리 예정

### Event-Sourcing

- [EVENT_SOURCING.md](./EVENT_SOURCING.md)에서 정리 예정

### CQRS

- [CQRS.md](./CQRS.md)에서 정리 예정

## 결론

결론부터 말하자면, 현재 상황에서 회사 서비스에 EDA 적용은 오버 엔지니어링이 아닙니다. 오히려 다음과 같은 이유로 적극적으로 적용해야 함을 깨달았습니다:

- 팬아웃 및 병렬 처리의 필요
    - 현재 회사 서비스에서 동일한 이벤트를 각 마이크로서비스에 알리기 위해 서로 다른 방식(Redis PubSub, REST API...)를 사용하고 있습니다.
    -  하지만 EDA를 적용한다면, 별도의 사용자 지정 코드를 작성할 필요 없이 이벤트를 분산시킬 수 있습니다.
        - 특히, REST API의 구조가 바뀌거나 PubSub 메시지의 구조가 바뀔 때마다 각 마이크로서비스를 수정해야 했던 번거로움을 없앨 수 있습니다.
- 이종 시스템의 통합
    - 위에서 말했듯이 회사 서비스는 다양한 마이크로서비스 간에 통신을 합니다.
    - 이벤트 브로커를 통해 마이크로서비스 간 간접 연결 및 상호 운용성을 구축하여 유연한 통신이 가능하게 할 수 있습니다.
- 확장성
    - MSA로의 전환이 끝난 시점이 아니라 전환 중이기 때문에 새로운 마이크로서비스가 계속 추가될 예정입니다.
    - 그렇기 때문에 전환이 끝난 시점이 아니라 전환 중인 지금이 오히려 더 좋은 시점이라 생각합니다.
    - 또한, EDA 적용을 통해 얻을 수 있는 느슨한 결합은 MSA 전환을 위해 꼭 충족시켜야 할 요구사항입니다.

이 밖에도 다양한 EDA 적용 사례 및 논문을 검토함으로써 EDA 적용이 필요한 또 다른 이유들도 확인할 수 있었습니다:

- 생산성 향상
    - 다음 링크들에서는 EDA 적용 후, 실제 실무에서 생산성이 향상함을 확인할 수 있었습니다.
        - https://www.rst.software/blog/event-driven-architecture?utm_source=chatgpt.com
        - https://journalwjaets.com/sites/default/files/fulltext_pdf/WJAETS-2025-0280.pdf
- 가용성 확보
    - Netflix, Uber를 비롯한 빅테크들의 서비스들은 EDA를 적용함으로써 거대한 트래픽에 대한 가용성을 보장했습니다.
        - 물론, 위의 빅테크들의 서비스들을 회사 서비스와 비교하는 것은 말이 안 되긴 하지만, 성공 사례가 있다는 것은 좋은 best practices가 될 수 있음을 의미한다 생각합니다.

물론 회사 서비스에 EDA를 적용함에 제약 사항이 없는 것은 아닙니다. 특히, 공부할 수록 이벤트 브로커로 어떤 스택을 사용해야 할 지 고민이 되었습니다. 개인적인 욕심으로는 Kafka를 사용해 보고 싶었지만, 제 욕심으로 Kafka를 무턱대고 사용하기에는 운영이 걱정되었습니다 (이거야 말로 오버 엔지니어링인 것 같습니다). 그래서 대안을 찾던 중 Redis Streams에 대해 알게 되었습니다. 결론적으로는 초기에는 Redis Streams를 사용해 EDA를 적용하고, 점차 Kafka로 마이그레이션 하는 방향으로 회사에 제안하려 합니다. 
(Redis Streams와 Kafka 비교 분석은 [REDIS_STREAMS_VS_KAFKA.md](./REDIS_STREAMS_VS_KAFKA.md)에 작성할 예정입니다.)

## References

- https://www.geeksforgeeks.org/system-design/event-driven-architecture-system-design/
- https://aws.amazon.com/what-is/eda/
- https://en.wikipedia.org/wiki/Event-driven_architecture
- https://f-lab.kr/insight/understanding-event-driven-architecture?gad_source=1&gad_campaignid=22368870602&gbraid=0AAAAACGgUFfWdvOrbX1LX33e6dZ-3dmsc&gclid=CjwKCAiAu67KBhAkEiwAY0jAlVn_GUc8t9R-aXwNTFTowdxkHNS-ESPkcVt5s__jS2WdMELbE9f72xoCc-MQAvD_BwE
- https://www.confluent.io/learn/event-driven-architecture/#how-it-works
- https://martinfowler.com/articles/201701-event-driven.html
- https://www.rst.software/blog/event-driven-architecture?utm_source=chatgpt.com
- https://journalwjaets.com/sites/default/files/fulltext_pdf/WJAETS-2025-0280.pdf
- https://icepanel.io/blog/2024-11-26-state-of-software-architecture-2024