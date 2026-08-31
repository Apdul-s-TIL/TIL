# Deployments

## Deployment란?

- 디플로이먼트(Deployment)는 보통 상태를 유지하지 않는(Stateless) 애플리케이션 워크로드를 실행하기 위해 파드 집합을 관리하는 리소스임
- 디플로이먼트는 파드와 레플리카셋(ReplicaSet)에 대한 선언적 업데이트(Declarative Updates)를 제공함
- 사용자가 디플로이먼트에 '의도한 상태(Desired State)'를 기술하면, 디플로이먼트 컨트롤러(Deployment Controller)가 제어된 속도(Controlled rate)로 현재 상태를 의도하는 상태로 변경함
- 새 리플리카셋을 생성하는 디플로이먼트를 정의하거나, 기존 디플로이먼트를 제거하고 모든 리소스를 새 디플로이먼트로 적용할 수 있음
  - 주의사항: 디플로이먼트가 소유한 레플리카셋을 직접 수정하거나 관리하면 안 됨

## Use Cases

- 다음은 디플로이먼트의 일반적인 유스케이스임

### Rollout

- 디플로이먼트를 생성하여 레플리카셋을 롤아웃함
- 레플리카셋은 백그라운드에서 파드를 생성함
- 롤아웃 상태를 확인하여 성공 여부를 파악할 수 있음

### 새로운 상태 선언

- 디플로이먼트의 `PodTemplateSpec`을 업데이트하여 파드의 새로운 상태를 선언함
- 새 레플리카셋이 생성되고, 디플로이먼트는 기존 레플리카셋을 축소(Scale down)하면서 새 레플리카셋을 점진적으로 확장(Scale up)하여 파드가 제어된 속도로 교체되도록 함
  - 이때 새 레플리카셋이 생성될 때마다 디플로이먼트의 수정 버전(Revision)이 업데이트 됨

### Rollback

- 현재 디플로이먼트 상태가 불안정한 경우 이전 리비전으로 롤백함
  - 롤백을 수행할 때도 리비전이 업데이트 됨

### Scaling

- 더 많은 로드를 위해 디플로이먼트를 스케일 업(확장)함

### Pause & Resume

- 디플로이먼트 롤아웃을 일시 중지한 후 `PodTemppateSpec`에 여러 수정 사항을 적용하고, 다시 재개하여 단 한 번의 새로운 롤아웃으로 반영함

### 롤아웃 정체 감지

- 디플로이먼트의 상태를 확인하여 롤아웃이 멈췄는지(Stuck) 파악함

### 오래된 레플리카셋 정리함

- 더 이상 필요 없는 이전 레플리카셋을 정리함

## 디플로이먼트 생성

### 예시 및 설명 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

- 위 예시는 3개의 `nginx` 파드를 실행하는 레플리카셋을 생성하는 디플로이먼트 예시임

#### 예시 상세 설명
- `.metadata.name`
  - `nginx-deployment`라는 이름의 디플로이먼트를 생성함
  - 이 이름은 향후 생성될 레플리카셋과 파드 이름의 기반이 됨
- `.spec.replicas`
  - 3개의 복제된 파드를 생성하도록 지정함
- `.spec.selector`
  - 생성된 레플리카셋이 어떤 파드를 관리할지 찾는 기준을 정의함
  - 위 예시에서는 파드 템플릿에 정의된 레이블(`app: nginx`)을 선택함
- `.spec.template`
  - `.metadata.labels`: 파드에 `app: nginx` 레이블을 부여함
  - `.spec.containers`: `nginx:1.14.2` 이미지를 사용하는 `nginx` 단일 컨테이너를 실행함

#### 생성 및 상태 확인 명령

```bash
# 디플로이먼트 생성
kubectl apply -f https://k8s.io/examples/controllers/nginx-deployment.yaml

# 진행 상태 확인
kubectl rollout status deployment/nginx-deployment

# kubectl rollout status deployment/nginx-deployment 실행 결과 예시
Waiting for rollout to finish: 2 out of 3 new replicas have been updated...
deployment "nginx-deployment" successfully rolled out

# 생성된 디플로이먼트 확인
kubectl get deployments

# kubectl get deployments 실행 결과 예시
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   0/3     0            0           1s

# 생성된 레플리카셋(RS) 확인
kubectl get rs

# kubectl get rs 실행 결과 예시
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-75675f5897   3         3         3       18s
```

- `kubectl get deployments` 실행 결과에 출력되는 주요 필드 의미는 다음과 같음:
  - NAME: 네임스페이스에 있는 디플로이먼트 이름의 목록
  - READY: 사용자가 사용할 수 있는 애플리케이션의 레플리카의 수를 표시함(`ready / desired` 패턴을 따름)
  - UP-TO-DATE: 의도한 상태를 얻기 위해 업데이트된 레플리카의 수를 표시함
  - AVAILABLE: 사용자가 사용할 수 있는 애플리케이션 레플리카의 수를 표시함
  - AGE: 애플리케이션의 실행된 시간을 표시함
- `kubectl get rs` 실행 결과에 출력되는 주요 필드 의미는 다음과 같음:
  - NAME: 네임스페이스에 있는 레플리카셋 이름의 목록
  - DESIRED: 디플로이먼트의 생성 시 정의된 의도한 애플리케이션 레플리카의 수를 표시함. 이것이 의도한 상태임
  - CURRENT: 현재 실행 중인 레플리카의 수를 표시함
  - READY: 사용자가 사용할 수 있는 애플리케이션의 레플리카의 수를 표시함
  - AGE: 애플리케이션의 실행된 시간을 표시함
- 레플리카셋의 이름 형식은 항상 `[DEPLOYMENT-NAME]-[HASH]` 형태임
  - 예: `nginx-deployment-75675f5897`
  - 생성된 `HASH` 값은 레플리카셋의 `pod-template-hash` 레이블과 동일함
- 주의사항
  - 다른 컨트롤러(다른 Deployment, StatefulSet 등)와 레이블/셀렉터가 겹치지 않도록 주의해야 함
  - 쿠버네티스는 중복 생성을 막지 않으므로, 셀렉터가 겹치면 컨트롤러 간 충돌이 발생해 예기치 못한 장애가 일어날 수 있음

### pod-template-hash 레이블

- 디플로이먼트 컨트롤러는 자신이 생성하거나 채택한 모든 레플리카셋에 `pod-template-hash` 레이블을 자동으로 추가함
- 이 레이블은 디플로이먼트의 자식 레플리카셋이 서로 중복되지 않도록 보장하는 역할을 함
  - 레플리카셋의 `PodTemplate`을 해싱함 
  - 해시 결과를 레플리카셋 셀렉터, 파드 템플릿 레이블 및 레플리카셋이 가질 수 있는 기존의 모든 파드에 레이블 값으로 추가해서 사용하도록 생성함
- 사용자가 직접 변경하면 안 됨

## 디플로이먼트 업데이트

- 디플로이먼트의 롤아웃은 파드 템플릿(즉, `.spec.template`)이 변경될 때만 트리거 됨
  - 예를 들면 템플릿의 레이블이나 컨테이너 이미지가 업데이트된 경우임
- 디플로이먼트의 스케일링 작ㅇ버은 롤아웃을 트리거하지 않음

### 이미지 업데이트 방법 및 예시

```bash
# CLI 명령어로 직접 업데이트
kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1

# 또는 매니페스트 직접 편집
kubectl edit deployment/nginx-deployment
```

- 업데이트가 진행되는 동안 `kubectl get rs`를 실행해 보면, 이전 레플리카셋의 파드 수는 줄어들고(Scale down to 0), 새 레플리카셋의 파드 수는 늘어나는(Scale up to 3) 과정을 확인할 수 있음

#### 롤업 업데이트 규칙(Max Unavailable & Max Surge)

- 기본적으로 디플로이먼트는 업데이트 중에 애플리케이션의 가용성을 보장함
- Max Unavailable(최대 불가 개수)
  - 업데이트 중 가용하지 않아도 되는 파드의 최대 비율/개수임
  - 기본값은 25%(의도한 파드 수의 최소 75%는 항상 유지됨)
- Max Surge(최대 초과 개수)
  - 의도한 파드 수를 초과하여 생성될 수 있는 최대 파드 비율/개수임
  - 기본값은 25%(의도한 파드 수의 최대 125%까지 순간 생성 가능)

### 롤오버(Rollover, 일명 인-플라이트 다중 업데이트)

- 디플로이먼트 컨트롤러는 각 시간마다 새로운 디플로이먼트에서 레플리카셋이 의도한 파드를 생성하고 띄우는 것을 주시함
  - 만약 디플로이먼트가 업데이트되면, 기존 레플리카셋에서 `.spec.selector` 레이블과 일치하는 파드를 컨트롤 하지만, 템플릿과 `.spec.template`이 불일치하면 스케일 다운이 됨
  - 결국 새로운 레플리카셋은 `.spec.replicas`로 스케일되고, 모든 기존 레플리카셋은 0개로 스케일 됨
- 만약 기존 롤아웃이 진행되는 중에 디플로이먼트를 업데이트하는 경우, 디플로이먼트가 업데이트에 따라 새 레플리카셋을 생성하고, 스케일 업하기 시작함
  - 그리고 이전에 스케일 업 하던 레플리카셋에 롤오버 함
  - 즉 기존에 생성 중이던 레플리카셋을 '이전 레플리카셋' 목록으로 넘기고 즉시 최신 변경 사항에 맞는 새로운 레플리카셋을 생성하여 전환함
- 예시
  - 디플로이먼트로 `nginx:1.14.2` 레플리카를 5개 생성을 해 `nginx:1.14.2` 레플리카가 3개 생성되었을 때 디플로이먼트를 업데이트해서 `nginx:1.16.1` 레플리카 5개를 생성하도록 업데이트를 한다고 가정
  - 이 경우 디플로이먼트는 즉시 생성된 3개의 `nginx:1.14.2` 파드 3개를 죽이기 시작하고 `nginx:1.16.1` 파드를 생성하기 시작함
  - 롤오버는 `nginx:1.14.2` 레플리카 5개가 생성되는 것을 기다리지 않음

### 레이블 셀렉터 업데이트

- 일반적으로 레이블 셀렉터(`.spec.selector`)는 변경 불가능(Immutable) 함
- 변경이 꼭 필요한 경우 디플로이먼트를 삭제 후 다시 생성해야 함
- 셀력터 추가 시 디플로이먼트의 사양에 있는 파드 템플릿 레이블도 새 레이블로 업데이트해야 함
  - 그렇지 않으면 유효성 검사 오류가 반환됨
  - 이 변경은 겹치지 않는 변경으로 새 셀렉터가 이전 셀렉터로 만든 레플리카셋과 파드를 선택하지 않게 됨
  - 그 결과 모든 기존 레플리카셋은 고아가 되며, 새로운 레플리카셋을 생성하게 됨
- 셀렉터 업데이트는 기존 셀렉터 키 값을 변경하며, 결과적으로 추가와 동일한 동작을 함
- 셀렉터 삭제는 디플로이먼트 셀렉터의 기존 키를 삭제하며 파드 템플릿 레이블의 변경을 필요로 하지 않음
  - 기존 레플리카셋은 고아가 아니고, 새 레플리카셋은 생성되지 않음
  - 그러나 제거된 레이블은 기존 파드와 레플리카셋에 여전히 존재함

## 디플로이먼트 롤백(Rolling Back)

- 기본적으로 모든 디플로이먼트의 롤아웃 기록은 시스템에 남아있어 언제든지 원할 때 롤백이 가능함
  - 예를 들어 디플로이먼트가 지속적인 충돌로 안정적이지 않은 경우, 롤백이 가능함
- 참고로 디플로이먼트의 수정 버전은 디플로이먼트 롤아웃시 생성됨
  - 이는 디플로이먼트 파드 템플릿(`.spec.template`)이 변경되는 경우에만 새로운 수정 버전이 생성된다는 것을 의미함
    - 예를 들어 템플릿의 레이블 또는 컨테이너 이미지를 업데이트 하는 경우가 있음
  - 디플로이먼트의 스케일링과 같은 다른 업데이트시 디플로이먼트 수정 버전은 생성되지 않으며, 수동-스케일링 또는 자동-스케일링을 동시에 수행할 수 있음
    - 이는 이전 수정 버전으로 롤백을 하는 경우에 디플로이먼트 파드 템플릿 부분만 롤백된다는 것을 의미함

### 예시

```
# 롤아웃 이력(Revision History) 확인
kubectl rollout history deployment/nginx-deployment

# 특정 리비전의 상세 정보 확인
kubectl rollout history deployment/nginx-deployment --revision=2

# 바로 이전 리비전으로 롤백
kubectl rollout undo deployment/nginx-deployment

# 특정 리비전으로 롤백
kubectl rollout undo deployment/nginx-deployment --to-revision=2
```

- 잘못된 이미지 이름(예: `nginx:1.161` 오타)을 적용하여 파드가 `ImagePullBackOff` 에러에 빠진 경우 이전 상태로 롤백하는 명령어임

## 디플로이먼트 스케일링(Scaling)

### 예시 및 설명

```
# 수동 스케일링 (파드 개수 조정)
kubectl scale deployment/nginx-deployment --replicas=10

# 오토스케일링 설정 (HPA 사용 - CPU 사용률 80% 기준)
kubectl autoscale deployment/nginx-deployment --min=10 --max=15 --cpu-percent=80%
```

### 비례적 스케일링(Proportional Scaling)

- 사용자 또는 오토스케일러가 롤아웃 중에 있는 디플로이먼트 롤링 업데이트를 스케일링 하는 경우(진행중 또는 일시 중지 중), 디플로이먼트 컨트롤러는 위험을 줄이기 위해 기존의 활성화된 여러 레플리카셋(파드와 레플리카셋)에 비례하여 파드 수를 분배해 추가함
  - 이를 proportional scailing 이라 함

#### 예시

```
# 1. 디플로이먼트에 있는 10개의 레플리카가 실행되는지 확인
kubectl get deploy

# 1번 실행 결과 예시
NAME                 DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment     10        10        10           10          50s

# 2. 클러스터 내부에서 확인할 수 없는 새 이미지로 업데이트 
kubectl set image deployment/nginx-deployment nginx=nginx:sometag

# 2번 실행 결과 예시
deployment.apps/nginx-deployment image updated

# 3. 롤아웃 상태 확인
kubectl get rs

# 3번 실행 결과 예시
NAME                          DESIRED   CURRENT   READY     AGE
nginx-deployment-1989198191   5         5         0         9s
nginx-deployment-618515232    8         8         8         1m

# 4. 디플로이먼트 확인 
kubectl get deploy

# 4번 실행 결과 예시
NAME                 DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment     15        18        7            8           7m

# 5. 롤아웃 상태 확인
kubectl get rs

# 5번 실행 결과 예시
NAME                          DESIRED   CURRENT   READY     AGE
nginx-deployment-1989198191   7         7         0         7m
nginx-deployment-618515232    11        11        11        7m
```

- 10개의 레플리카를 디플로이먼트로 `maxSurge=3`, 그리고 `maxUnavailable=2`로 실행한다고 가정한 시나리오임
- `3. 롤아웃 상태 확인`에서 `nginx-deployment-1989198191`(예시)로 새로운 롤 아웃이 시작되지만, 위에서 언급한 `maxUnavailable`의 요구 사항으로 차단 됨
  - 그 다음 디플로이먼트에 대한 새로운 스케일링 요청이 따라 옴
  - 오토스케일러는 디플로이먼트 레플리카를 15로 증가시킴
  - 디플로이먼트 컨트롤러는 새로운 4개의 레플리카의 추가를 위한 위치를 결정해야 함
  - 만약 비례적 스케일링을 사용하지 않으면 4개 모두 새 레플리카셋에 추가됨
  - 비례적 스케일링으로 추가 레플리카를 모든 모든 레플리카셋에 걸쳐 분산할 수 있음
  - 비율이 높을수록 가장 많은 레플리카가 있는 레플리카셋으로 이동하고, 비율이 낮을 수록 적은 레플리카가 있는 레플리카셋으로 이동함
  - 남은 것들은 대부분의 레플리카가 있는 레플리카셋에 추가됨
  - 0개의 레플리카가 있는 레플리카셋은 스케일 업 되지 않음
- 위 예시에서 기존 레플리카셋에 3개의 레플리카가 추가되고, 2개의 레플리카는 새 레플리카에 추가됨
  - 결국 롤아웃 프로세스는 새 레플리카가 정상이라고 가정하면 모든 레플리카를 새 레플리카셋으로 이동시킴

## 디플로이먼트 롤아웃 일시 중지와 재개(Pausing and Resuming)

- 여러 변경 사항에(이미지 변경, 리소스 제한 변경 등)을 한 번에 적용하면서 불필요한 중간 롤아웃을 방지하고 싶을 때 사용함

### 예시

```
# 1. 롤아웃 일시 중지
kubectl rollout pause deployment/nginx-deployment

# 2. 원하는 여러 변경 사항 적용
kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1
kubectl set resources deployment/nginx-deployment -c=nginx --limits=cpu=200m,memory=512Mi

# (이 동안에는 새 롤아웃이 시작되지 않음)

# 3. 롤아웃 재개
kubectl rollout resume deployment/nginx-deployment
```

- 일시 중지된 디플로이먼트를 재개할 때까지 롤백(`rollout undo`)할 수 없음

## 디플로이먼트 상태(Deployment Status)

- 디플로이먼트는 생명주기(라이프사이클) 동안 다음과 같은 상태를 거침
  - 진행 중(Progressing)
    - 새 레플리카셋 생성, 스케일 업/다운, 새 파드의 준비 상태 확인 등이 진행 중인 상태
    - `type: Progressing`, `status: "True"`
  - 완료됨(Complete)
    - 요청한 모든 복제본이 최신 버전으로 업데이트되었고, 사용 가능하며, 이전 파드가 더 이상 실행되지 않는 상태
    - `reason: NewReplicaSetAvailable`
  - 진행 실패(Failed)
    - 쿼터 부족, Readiness Probe 실패, 이미지 풀 에러, 권한 부족 등으로 진행이 중단된 상태
    - `.spec.progressDeadlineSeconds`(기본값: 600초/10분) 시간을 초과하면 `reason: ProgressDeadlineExceeded` 에러 조건이 기록됨

### 디플로이먼트 진행 중(Progressing)

- 쿠버네티스는 다음 작업 중 하나를 수행할 때 디플로이먼트를 `Progressing`으로 표시함
  - 디플로이먼트로 새 레플리카셋을 생성
  - 디플로이먼트로 새로운 레플리카셋을 스케일 업
  - 디플로이먼트로 기존 레플리카셋을 스케일 다운
  - 새 파드가 준비되거나 이용할 수 있음(최소 준비 시간(초) 동안 준비됨)
- 롤아웃이 `Progressing` 상태가 되면, 디플로이먼트 컨트롤러는 디플로이먼트의 `.status.conditions`에 다음 속성을 포함하는 컨디션을 추가함
  - `type: Progressing`
  - `status: "True"`
  - `reason: NewReplicaSetCreated` | `reason: FoundNewReplicaSet` | `reason: ReplicaSetUpdated`
- `kubectl rollout status`를 사용해서 디플로이먼트의 진행상황을 모니터할 수 있음

### 디플로이먼트 완료(Complete)

- 쿠버네티스는 다음과 같은 특성을 가지게 되면 디플로이먼트를 `Complete`로 표시함
  - 디플로이먼틀와 관련된 모든 레플리카가 지정된 최신 버전으로 업데이트 되었을 때
    - 즉, 요청한 모든 업데이트가 완료되었을 때
  - 디플로이먼트와 관련한 모든 레플리카를 사용할 수 있을 때
  - 디플로이먼트에 대해 이전 복제본이 실행되고 있지 않을 때
- 롤아웃이 `Complete` 상태가 되면, 디플로이먼트 컨트롤러는 디플로이먼트의 `.status.conditions`에 다음 속성을 포함하는 컨디션을 추가함
  - `type: Progressing`
  - `status: "True"`
  - `reason: NewReplicaSetAvailable`
  - 이 `Progressing` 컨디션은 새로운 롤아웃이 시작되기 전까지는 `"True"` 상태값을 유지할 것임
    - 레플리카의 가용성이 변경되는 경우에도(이 경우 `Available` 컨디션에 영향을 미침) 컨디션은 유지됨
- `kubectl rollout status`를 사용해서 디플로이먼트가 완료되었는지 확인할 수 있음
  - 만약 롤아웃이 성공적으로 완료되면, `kubectl rollout status` 는 종료 코드로 0이 반환됨

### 디플로이먼트 실패(Failed)

- 디플로이먼트가 새로운 레플리카셋을 배포하는 과정에서 완료되지 못한 채 중단(Stuck)될 수 있고, 이러한 상황은 주로 다음과 같은 원인으로 발생함
  - 리소스 쿼터(Quota, 할당량) 부족
  - Readiness probe의 실패
  - 이미지 풀(Image pull) 에러
  - 권한 부족
  - 범위 제한(Limit ranges) 제약
  - 애플리케이션 런타임의 잘못된 구성
- 이러한 배포 중단 상태를 감지하는 방법 중 하나는 디플로이먼트 명세에 데드라인 파라미터(`.spec.progressDeadineSeconds`)를 지정하는 것임
  - `.spec.progressDeadlineSeconds`는 디플로이먼트 컨트롤러가 롤아웃 진행이 정체되었음을 디플로이먼트 상태(Status)에 표시하기까지 대기하는 시간(초 단위)을 의미함
  - 다음과 같은 `kubectl` 명령어로 `progressDeadlineSeconds`를 설정해 컨트롤러가 10분 후 디플로이먼트 롤아웃에 대한 진행 상태의 부족에 대한 리포트를 수행하게 할 수 있음
    - `kubectl patch deployment/nginx-deployment -p '{"spec":{"progressDeadlineSeconds":600}}'`
- 지정한 데드라인 시간을 초과하면, 디플로이먼트 컨트롤러는 `.status.conditions` 속성에 다음과 같은 `DeploymentCondition`을 추가함
  - `type: Progressing`
  - `status: "False"`
  - `reason: ProgressDeadlineExceeded`
  - 이 컨디션은 `ReplicaSetCreateError`와 같은 이유로 데드라인 시간 이전이라도 조기에 실패하여 `status` 값이 `"False"`로 설정될 수 있음
  - 또한, 디플로이먼트 롤아웃이 성공적으로 완료되고나면 이 데드라인은 더 이상 고려되지 않음

#### 참고

- 쿠버네티스는 `reason: ProgressDeadlineExceeded`과 같은 상태 조건을 보고하는 것 외에 정지된 디플로이먼트에 대해 조치를 취하지 않음
- 더 높은 수준의 오케스트레이터는 이를 활용할 수 있음
  - 예를 들어 디플로이먼트를 이전 버전으로 롤백할 수 있음
- 만약 디플로이먼트 롤아웃을 일시 중지하면 쿠버네티스는 지정된 데드라인과 비교하여 진행 상황을 확인하지 않음
  - 롤아웃 중에 디플로이먼트 롤아웃을 안전하게 일시 중지하고, 데드라인을 넘기도록 하는 조건을 트리거하지 않고 재개할 수 있음 

### 실패한 디플로이먼트에서의 운영

- 완료된 디플로이먼트에 적용되는 모든 행동은 실패한 디플로이먼트에도 적용됨
- 디플로이먼트 파드 템플릿에서 여러 개의 수정사항을 적용해야 하는 경우, 다음과 같은 방식으로 운영할 수 있음
  - 스케일 업/다운
  - 이전 수정 버전으로 롤백
  - 일시 중지

## 정책 초기화

- 디플로이먼트의 `.spec.revisionHistoryLimit` 필드를 설정해서 디플로이먼트에서 유지해야 하는 이전 레플리카셋의 수를 명시할 수 있음
  - 나머지는 백그라운드에서 가비지 컬렉션(Garbage collection) 처리됨
  - 기본값은 10임
- 참고사항
  - 이 필드를 0으로 설정하면 디플로이먼트의 기록을 전부 초기화 함
  - 이로 인해 디플로이먼트를 이전 상태로 롤백할 수 없게 됨
- 초기화 작업은 디플로이먼트가 `완료(Completion)` 상태에 도달한 이후에 시작됨
  - `.spec.revisionHistoryLimit`을 0으로 설정하더라도 쿠버네티스가 이전 레플리카셋을 삭제하기 전에 반드시 새로운 레플리카셋을 먼저 생성함
- `.spec.revisionHistoryLimit` 값을 0이 아닌 값으로 설정하더라도, 실제 레플리카셋 개수가 설정한 제한을 초과할 수 있음
  - 예를 들어, 파드가 무한 충돌(Crash looping) 상태에 빠지고, 시간에 따라 여러 번의 롤링 업데이트가 발생하면 디플로이먼트가 `완료(Completion)` 상태에 도달하지 못하므로 `.spec.revisionHistoryLimit`에 설정된 수보다 더 많은 레플리카셋이 남아있을 수 있음

## 카나리 디플로이먼트(Canary Deployment)

- 디플로이먼트를 사용하여 사용자나 서버의 일부 집합에만 새 릴리스를 롤아웃하려는 경우, 카나리 패턴에 따라 각 릴리스별로 여러 디플로이먼트를 생성할 수 있음

## 디플로이먼트 명세 작성하기(Writing a Deployment Spec)

- 다른 모든 쿠버네티스 설정과 마찬가지로 디플로이먼트에도 다음과 같은 필드가 필요함
  - `.apiVersion`
  - `.kind`
  - `.metadata`
- 컨트롤 플레인이 디플로이먼트를 위한 새 파드를 생성할 때, 디플로이먼트의 `metadata.name`은 해당 파드들의 이름을 짓는 기준의 일부가 됨
  - 디플로이먼트 오브젝트의 이름은 유효한 DNS 서브도메인 값이어야 함
  - 그렇지 않으면 파드 호스트네임에 예상치 못한 결과를 초래할 수 있음
  - 호환성을 극대화하려면 좀 더 엄격한 규칙인 DNS 레이블 규칙을 따르는 이름이 좋음
- 디플로이먼트에는 `.spec` 섹션도 필요함

### 파드 템플릿(Pod Template)

- `.spec.template`와 `.spec.selector`는 `.spec` 섹션에서 유일한 필수 필드임
- `.spec.template`는 파드 템플릿임
  - 중첩 구조이며 `apiVersion`이나 `kind`가 없다는 점을 제외하면 파드와 동일한 스키마를 가지고 있음
- 디플로이먼트 내의 파드 템플릿은 파드의 필수 필드 외에도 적절한 레이블과 적절한 재시작 정책(Restart policy)을 명시해야 함
- 레이블의 경우 다른 컨트롤러와 중복되지 않도록 주의해야 함
- `.spec.template.spec.restartPolicy`는 오직 `Always`만 허용되며, 명시하지 않을 경우 기본값으로 지정됨

### 레플리카(Replicas)

- `.spec.replicas`는 원하는 파드의 개수를 지정하는 선택적 필드임
- 기본값은 1임
- 예를 들어 `kubectl scale deployment deployment --replicas=X` 명령을 통해 디플로이먼트를 수동으로 스케일링한 후, 매니페스트 기반(예: `kubectl apply -f deployment.yaml` 실행)으로 디플로이먼트를 업데이트하면, 수동으로 설정했던 디플로이먼트의 크기가 오버라이드 됨
- `HorizontalPodAutoscaler`(또는 수평 스케일링을 위한 유사한 API)가 디플로이먼트 크기를 관리하고 있다면, `.spec.replicas`를 설정하면 안 됨
  - 대신 쿠버네티스 컨트롤 플레인이 `.spec.replicas` 필드를 자동으로 관리하도록 두는 것이 좋음

### 셀렉터(Selector)

- `.spec.selector`는 디플로이먼트가 관리할 대상 파드를 지정하는 레이블 셀렉터이며 필수 필드임
- `.spec.selector`는 `.spec.template.metadata.labels`와 반드시 일치해야 하며, 그렇지 않으면 API에 의해 거부됨
- `apps/v1` API 버전에서는 `.spec.selector`와 `.metadata.labels`가 설정되지 않았더라도 `.spec.template.metadata.labels` 값으로 자동 설정되지 않았음
  - 따라서 명시적으로 설정해야 함
  - 또한, `apps/v1`에서 `.spec.selector`는 디플로이먼트 생성 후 변경할 수 없음(Immutable)
- 디플로이먼트는 템플릿이 `.spec.template`과 다르거나 파드의 총개수가 `.spec.replicas`를 초과하는 경우, 셀렉터와 레이블이 일치하는 파드를 종료할 수 있음
  - 반대로 파드 수가 원하는 수보다 적으면 `.spec.template`를 바탕으로 새 파드를 생성함
- 참고사항
  - 다른 디플로이먼트를 생성하거나 레플리카셋 또는 레플리케이션 컨트롤러(ReplicationController) 같은 다른 컨트롤러를 사용해서 직접적으로 레이블과 셀렉터가 일치하는 다른 파드를 생성하지 말아야 함
    - 그렇게 하면 첫 번째 디플로이먼트는 그 파드들마저 자기가 생성한 것으로 간주함
    - 쿠버네티스는 이를 금지하지 않으므로 주의가 필요함
    - 중복된 셀렉터를 가진 컨트롤러가 여러 개 존재하면, 컨트롤러들이 서로 충돌하여 올바르게 작동하지 않음

### 전략(Strategy)

- `.spec.strategy`는 이전 파드를 새로운 파드로 교체할 때 사용하는 전략을 지정함
- `.spec.strategy.type`은 `Recreate`(재생성) 또는 `RollingUpdate`(롤링업데이트)가 될 수 있음
  - 기본값은 `RollingUpdate`임

#### Recreate 디플로이먼트

- `.spec.strategy.type==Recreate`로 설정하면 새 파드가 생성되기 전에 기존의 모든 파드가 먼저 종료됨
- 참고사항
  - 이 전략은 업그레이드 시 생성 직전에 파드가 종료되는 것만 보장함
  - 디플로이먼트를 업그레이드하면 이전 버전의 모든 파드가 즉시 종료됨
  - 신규 버전의 파드가 생성되기 전에 성공적으로 제거가 완료되기를 대기함
  - 파드를 수동으로 삭제하면, 라이프사이클은 레플리카셋에 의해 제어되므로(이전 파드가 여전히 `Terminating` 상태이더라도) 대체 파드가 즉시 재성성 됨
  - 파드에 대해 "최대" 보장이 필요한 경우 StatefulSet 사용을 고려해야 함

#### RollingUpdate 디플로이먼트

- `.spec.strategy.type==RollingUpdate`로 설정하면 디플로이먼트가 롤링 업데이트 방식으로 파드를 업데이트함
- `maxUnavailable`과 `maxSurge`를 지정하여 롤링 업데이트 프로세스를 제어할 수 있음
- Max Unavailable(최대 불가)
  - `.spec.strategy.rollingUpdate.maxUnavailable`은 업데이트 프로세스 중 사용할 수 없는 파드의 최대 개수를 지정하는 선택적 필드임
  - 절대적인 수치(e.g. 5) 또는 원하는 파드 비율(e.g. 10%)로 설정할 수 있음
    - 비율 설정 시, 내림하여 계산함
  - `.spec.strategy.rollingUpdate.maxSurge`가 0인 경우, 이 값은 0이 될 수 없음
  - 기본값은 25%임
  - 예를 들어 이 값을 30%로 설정하면 롤링업데이트 시작시, 즉각 이전 레플리카셋의 크기를 의도한 파드 중 70%로 스케일 다운할 수 있음
    - 새 파드가 준비되면 기존 레플리카셋을 스케일 다운할 수 있으며, 업데이트 중에 항상 사용 가능한 전체 파드의 수는 의도한 파드의 수의 70% 이상이 되도록 새 레플리카셋을 스케일 업 할 수 있음
- Max Surge(최대 초과 개수)
  - `.spec.strategy.rollingUpdate.maxSurge`는 의도한 파드 수를 초과하여 생성할 수 있는 최대 파드 개수를 지정하는 선택적 필드임
  - 절대적인 수치(e.g. 5) 또는 원하는 파드 비율(e.g. 10%)로 설정할 수 있음
    - 비율 설정 시, 올림하여 계산함
  - `.spec.strategy.rollingUpdate.maxUnavailable`이 0인 경우, 이 값은 0이 될 수 없음
  - 기본값은 25%임
  - 예를 들어 이 값을 30%로 설정하면 롤링업데이트 시작시, 새 레플리카셋의 크기를 즉시 조정해 기존 및 새 파드의 전체 갯수를 의도한 파드의 130%를 넘지 않도록 함
    - 기존 파드가 죽으면 새로운 레플리카셋은 스케일 업할 수 있으며, 업데이트하는 동안 항상 실행하는 총 파드의 수는 최대 의도한 파드의 수의 130%r가 되도록 보장함

### 진행 데드라인 시간(Progress Deadline Seconds)

- `.spec.progressDeadlineSeconds`는 디플로이먼트가 표면적으로 `type: Progressing`, `status: "False"`의 상태 그리고 리소스가 `reason: ProgressDeadlineExceeded` 상태로 진행 실패를 보고하기 전에 디플로이먼트가 진행되는 것을 대기시키는 시간(초)를 명시하는 선택적 필드임
- 디플로이먼트 컨트롤러는 디플로이먼트를 계속 재시작 함
- 기본값은 600초임
- 자동화된 롤백이 구현된다면 디플로이먼트 컨트롤러는 상태를 관찰하고, 그 즉시 디플로이먼트를 롤백함
- 만약 명시된다면 이 필드는 `.spec.minReadySeconds`보다 커야 함

### 최소 준비 시간(Min Ready Seconds)

- `.spec.minReadySeconds`는 새로 생성된 파드가 사용 가능한(Available) 상태로 간주되기 위해 컨테이너 크래시 없이 준비 상태를 유지해야 하는 최소 시간(초)을 지정하는 선택적 필드임
- 기본값은 0(준비되는 즉시 사용 가능으로 간주)임

### 종료 중인 파드(Terminating Pods)

- API 서버와 `kube-controller-manager`에서 `DeploymentReplicaSetTerminatingReplicas` 기능 게이트가 활성화되어 있어야 종료 중인 파드를 확인할 수 있음
- 삭제나 스케일 다운으로 인해 종료되는 파드는 정리되기까지 오랜 시간이 걸릴 수 있으며, 그동안 추가 리소스를 소비할 수 있음
  - 이로 인해 전체 파드 수가 일시적으로 `.spec.replicas`를 초과할 수 있음
- 종료 중인 파드는 디플로이먼트의 `.status.terminatingReplicas` 필드를 통해 추적할 수 있음

### 수정 버전 기록 제한을(Revision History Limit)

- 디플로이먼트의 수정 버전(리비전) 기록은 자신이 컨트롤하는 레플리카셋에 저장됨
- `.spec.revisionHistoryLimit`은 롤백을 허용하기 위해 보관할 이전 레플리카셋의 개수를 지정하는 선택적 필드임
- 이전 레플리카셋은 etcd 리소스를 차지하고 `kubectl get rs` 출력 결과를 복잡하게 만듬
- 각 디플로이먼트의 구성은 디플로이먼트의 레플리카셋에 저장됨
  - 이전 레플리카셋이 삭제되면 해당 디플로이먼트의 수정 버전으로 롤백할 수 있는 기능이 사라짐
  - 기본적으로 10개의 기존 레플리카셋이 유지되지만 이상적인 값은 새로운 디플로이먼트의 빈도와 안정성에 따라 달라짐
- 이 필드를 0으로 설정하면 레플리카가 0이 되며, 이전 레플리카셋이 정리됨
  - 이 경우, 새로운 디플로이먼트 롤아웃을 취소할 수 없음
  - 새로운 디플로이먼트의 롤아웃은 수정 버전 이력이 정리되기 때문임

### 일시 정지(Paused)

- `.spec.paused`는 디플로이먼트를 일시 중지하거나 재개하기 위한 선택적 불리언(Boolean) 필드임
- 일시 중지된 디플로이먼트는 일시 중지 상태가 유지되는 동안 `PodTemplateSpec` 변경 사항이 새로운 롤아웃을 트리거하지 읺음
- 기본값은 `false`임









