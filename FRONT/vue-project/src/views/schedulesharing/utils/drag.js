/**
 * 드래그 시작 이벤트 핸들러
 * @param {Event} event - 드래그 이벤트 객체
 * @param {string|number} memberId - 드래그되는 멤버의 ID
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function handleDragStart(event, memberId, context) {
  context.draggedMember = memberId;
  event.dataTransfer.setData('text/plain', memberId);
  event.dataTransfer.effectAllowed = 'move';
}

/**
 * 드래그 종료 이벤트 핸들러
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function handleDragEnd(context) {
  context.draggedMember = null;
}

/**
 * 드래그 오버 이벤트 핸들러
 * @param {Event} event - 드래그 오버 이벤트 객체
 */
export function handleDragOver(event) {
  event.preventDefault();
  event.dataTransfer.dropEffect = 'move';
}

/**
 * 드롭 이벤트 핸들러
 * @param {Event} event - 드롭 이벤트 객체
 * @param {string|number} targetId - 드롭 대상의 ID
 * @param {Function} callback - 드롭 후 실행할 콜백 함수
 */
export function handleDrop(event, targetId, callback) {
  event.preventDefault();
  const draggedId = event.dataTransfer.getData('text/plain');
  
  if (draggedId && draggedId !== targetId) {
    callback(draggedId, targetId);
  }
}

/**
 * 노드 위치 업데이트 (드롭 후 위치 계산)
 * @param {Event} event - 드롭 이벤트 객체
 * @param {string|number} draggedMemberId - 드래그된 멤버 ID
 * @param {Function} updateNodePosition - 노드 위치 업데이트 함수
 */
export function handleNodeDrop(event, draggedMemberId, updateNodePosition) {
  event.preventDefault();
  if (draggedMemberId) {
    // 드롭된 위치로 노드 이동
    const rect = event.currentTarget.getBoundingClientRect();
    const left = event.clientX - rect.left;
    const top = event.clientY - rect.top;
    updateNodePosition(draggedMemberId, left, top);
  }
}

/**
 * 노드 위치 업데이트 함수
 * @param {string|number} memberId - 멤버 ID
 * @param {number} left - 새로운 X 좌표
 * @param {number} top - 새로운 Y 좌표
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function updateNodePosition(memberId, left, top, context) {
  const member = context.teamMembers.find(m => m.id === memberId);
  const department = context.departments.find(d => d.name === member.department);
  
  if (department) {
    // 부서 영역 안에 드롭된 경우
    if (left >= department.x && left <= department.x + department.width &&
        top >= department.y && top <= department.y + department.height) {
      
      // 초기 위치로 자동 고정
      const initialPosition = getInitialNodePosition(memberId, context);
      context.nodePositions[memberId] = initialPosition;
      
      // 부드러운 애니메이션으로 초기 위치로 이동
      context.$nextTick(() => {
        animateToInitialPosition(memberId, initialPosition);
      });
      
      console.log(`${member.name}이(가) ${department.name} 영역 안에 드롭되어 초기 위치로 고정됩니다.`);
      
    } else {
      // 부서 영역 밖에 드롭된 경우 - nodePositions에서 제거하여 다음 드래그 시 초기 위치로 돌아가도록 함
      delete context.nodePositions[memberId];
      console.log(`${member.name}이(가) 자유 위치에 배치되었습니다. (다음 드래그 시 초기 위치로 돌아갑니다)`);
    }
  }
}

/**
 * 초기 위치 계산
 * @param {string|number} memberId - 멤버 ID
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @returns {Object} 초기 위치 객체
 */
function getInitialNodePosition(memberId, context) {
  const index = context.teamMembers.findIndex(m => m.id === memberId);
  return getNodePosition(index, context);
}

/**
 * 노드 위치 계산 (부서별 그룹화)
 * @param {number} index - 팀원 배열 인덱스
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @returns {Object} 위치 객체
 */
function getNodePosition(index, context) {
  const member = context.teamMembers[index];
  const department = member.department;
  
  // 부서별 위치 계산
  let baseX = 0;
  let baseY = 0;
  let offsetX = 0;
  let offsetY = 0;
  
  // 부서별 그룹 위치
  const departmentObj = context.departments.find(d => d.name === department);
  if (departmentObj) {
    // 부서 네모 중앙에 배치하되 아래로 내림
    baseX = departmentObj.x + departmentObj.width / 2 - 120;
    baseY = departmentObj.y + departmentObj.height / 2 - 50;
  }
  
  // 같은 부서 내에서 위치 조정
  const departmentMembers = context.teamMembers.filter(m => m.department === department);
  const memberIndexInDepartment = departmentMembers.findIndex(m => m.id === member.id);
  
  // 부서 내 행과 열 계산 (2x2 그리드 유지)
  const row = Math.floor(memberIndexInDepartment / 2);
  const col = memberIndexInDepartment % 2;
  
  // 간격을 크게 늘림
  offsetX = col * 240;
  offsetY = row * 220;
  
  const x = baseX + offsetX;
  const y = baseY + offsetY;
  
  return { left: `${x}px`, top: `${y}px` };
}

/**
 * 부드러운 애니메이션으로 초기 위치로 이동
 * @param {string|number} memberId - 멤버 ID
 * @param {Object} targetPosition - 목표 위치
 */
function animateToInitialPosition(memberId, targetPosition) {
  const nodeElement = document.querySelector(`[data-member-id="${memberId}"]`);
  if (nodeElement) {
    nodeElement.style.transition = 'all 0.5s ease-in-out';
    nodeElement.style.left = targetPosition.left;
    nodeElement.style.top = targetPosition.top;
    
    // 애니메이션 완료 후 transition 제거
    setTimeout(() => {
      nodeElement.style.transition = '';
    }, 500);
  }
}
