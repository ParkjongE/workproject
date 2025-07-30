<template>
  <Container>
    <SideBar />
    <MainContent>
      <MainTitle>조직도</MainTitle>
      <OrgChartArea id="orgchart-area">
        <OrgLevel id="level1">
          <OrgNode class="org-node" id="node-1">영업부</OrgNode>
          <OrgNode class="org-node" id="node-2">기술부</OrgNode>
          <OrgNode class="org-node" id="node-3">마케팅부</OrgNode>
          <OrgNode class="org-node" id="node-4">인사부</OrgNode>
          <OrgNode class="org-node" id="node-5">경영부</OrgNode>
        </OrgLevel>
        <OrgLevel id="level2">
          <OrgNode class="org-node">영업1팀</OrgNode>
          <OrgNode class="org-node">영업2팀</OrgNode>
          <OrgNode class="org-node">기술1팀</OrgNode>
          <OrgNode class="org-node">기술2팀</OrgNode>
          <OrgNode class="org-node">마케팅팀</OrgNode>
          <OrgNode class="org-node">인사팀</OrgNode>
          <OrgNode class="org-node">경영지원팀</OrgNode>
        </OrgLevel>
        <OrgLevel id="level3">
          <OrgNode class="org-node">A파트</OrgNode>
          <OrgNode class="org-node">B파트</OrgNode>
          <OrgNode class="org-node">C파트</OrgNode>
          <OrgNode class="org-node">D파트</OrgNode>
          <OrgNode class="org-node">E파트</OrgNode>
          <OrgNode class="org-node">F파트</OrgNode>
          <OrgNode class="org-node">G파트</OrgNode>
        </OrgLevel>
        <OrgSvg id="org-svg" ref="orgSvg"></OrgSvg>
      </OrgChartArea>
    </MainContent>
  </Container>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  Container, MainContent, MainTitle, OrgChartArea, OrgLevel, OrgNode, OrgSvg
} from './OrgChart.styled.js';

export default {
  name: 'OrgChart',
  components: {
    SideBar,
    Container, MainContent, MainTitle, OrgChartArea, OrgLevel, OrgNode, OrgSvg
  },
  mounted() {
    this.$nextTick(() => {
      this.drawConnections();
      window.addEventListener('resize', this.drawConnections);
    });
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.drawConnections);
  },
  methods: {
    drawConnections() {
      const svg = this.$refs.orgSvg;
      if (!svg) return;
      // SVG 내용 초기화
      while (svg.firstChild) {
        svg.removeChild(svg.firstChild);
      }
      // 레벨 1과 레벨 2 연결
      const level1Nodes = document.querySelectorAll('#level1 .org-node');
      const level2Nodes = document.querySelectorAll('#level2 .org-node');
      if (!level1Nodes.length || !level2Nodes.length) return;
      // 영업부 -> 영업1팀, 영업2팀
      if (level1Nodes[0] && level2Nodes[0]) this.drawLine(level1Nodes[0], level2Nodes[0], svg);
      if (level1Nodes[0] && level2Nodes[1]) this.drawLine(level1Nodes[0], level2Nodes[1], svg);
      // 기술부 -> 기술1팀, 기술2팀
      if (level1Nodes[1] && level2Nodes[2]) this.drawLine(level1Nodes[1], level2Nodes[2], svg);
      if (level1Nodes[1] && level2Nodes[3]) this.drawLine(level1Nodes[1], level2Nodes[3], svg);
      // 마케팅부 -> 마케팅팀
      if (level1Nodes[2] && level2Nodes[4]) this.drawLine(level1Nodes[2], level2Nodes[4], svg);
      // 인사부 -> 인사팀
      if (level1Nodes[3] && level2Nodes[5]) this.drawLine(level1Nodes[3], level2Nodes[5], svg);
      // 경영부 -> 경영지원팀
      if (level1Nodes[4] && level2Nodes[6]) this.drawLine(level1Nodes[4], level2Nodes[6], svg);
      // 레벨 2와 레벨 3 연결
      const level3Nodes = document.querySelectorAll('#level3 .org-node');
      if (!level3Nodes.length) return;
      // 영업1팀 -> A파트
      if (level2Nodes[0] && level3Nodes[0]) this.drawLine(level2Nodes[0], level3Nodes[0], svg);
      // 영업2팀 -> B파트
      if (level2Nodes[1] && level3Nodes[1]) this.drawLine(level2Nodes[1], level3Nodes[1], svg);
      // 기술1팀 -> C파트, D파트
      if (level2Nodes[2] && level3Nodes[2]) this.drawLine(level2Nodes[2], level3Nodes[2], svg);
      if (level2Nodes[2] && level3Nodes[3]) this.drawLine(level2Nodes[2], level3Nodes[3], svg);
      // 기술2팀 -> E파트
      if (level2Nodes[3] && level3Nodes[4]) this.drawLine(level2Nodes[3], level3Nodes[4], svg);
      // 마케팅팀 -> F파트
      if (level2Nodes[4] && level3Nodes[5]) this.drawLine(level2Nodes[4], level3Nodes[5], svg);
      // 인사팀 -> G파트
      if (level2Nodes[5] && level3Nodes[6]) this.drawLine(level2Nodes[5], level3Nodes[6], svg);
    },
    drawLine(fromNode, toNode, svg) {
      if (!fromNode || !toNode || !svg) return;
      if (!fromNode.getBoundingClientRect || !toNode.getBoundingClientRect || !svg.getBoundingClientRect) return;
      const fromRect = fromNode.getBoundingClientRect();
      const toRect = toNode.getBoundingClientRect();
      const svgRect = svg.getBoundingClientRect();
      const fromX = fromRect.left + fromRect.width / 2 - svgRect.left;
      const fromY = fromRect.bottom - svgRect.top;
      const toX = toRect.left + toRect.width / 2 - svgRect.left;
      const toY = toRect.top - svgRect.top;
      const line = document.createElementNS('http://www.w3.org/2000/svg', 'path');
      const midY = (fromY + toY) / 2;
      line.setAttribute('d', `M${fromX},${fromY} C${fromX},${midY} ${toX},${midY} ${toX},${toY}`);
      line.setAttribute('stroke', '#ccc');
      line.setAttribute('stroke-width', '1');
      line.setAttribute('fill', 'none');
      svg.appendChild(line);
    }
  }
};
</script>

<style scoped>
/* CSS는 별도 파일로 분리됨 */
</style> 