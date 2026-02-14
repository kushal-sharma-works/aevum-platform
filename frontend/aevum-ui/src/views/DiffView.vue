<script setup lang="ts">
import PageContainer from '@/components/layout/PageContainer.vue'
import DecisionDiff from '@/components/decisions/DecisionDiff.vue'

const left = { status: 'Evaluated', amount: 42 }
const right = { status: 'Skipped', amount: 43, reason: 'threshold' }
</script>

<template>
	<PageContainer title="Decision Diff" description="Compare outputs across time or rule version">
		<DecisionDiff :left="left" :right="right" />
	</PageContainer>
</template>
