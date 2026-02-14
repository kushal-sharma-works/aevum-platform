<script setup lang="ts">
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import type { Decision } from '@/types/decision'
import DecisionTrace from './DecisionTrace.vue'

defineProps<{ decision: Decision }>()
</script>

<template>
	<div class="q-gutter-md">
		<q-card flat bordered>
			<q-card-section>
				<div class="text-subtitle1 text-weight-medium">Decision {{ decision.decisionId }}</div>
				<div class="text-caption text-grey-7">Rule {{ decision.ruleId }} v{{ decision.ruleVersion }}</div>
			</q-card-section>
		</q-card>
		<div class="row q-col-gutter-md">
			<div class="col-12 col-lg-6">
				<BaseJsonViewer :value="decision.input" label="Input" />
			</div>
			<div class="col-12 col-lg-6">
				<BaseJsonViewer :value="decision.output" label="Output" />
			</div>
		</div>
		<DecisionTrace :trace="decision.trace" />
	</div>
</template>
