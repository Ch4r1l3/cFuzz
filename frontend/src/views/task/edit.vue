<template>
  <div class="app-container">
    <task-detail :is-edit="true" />
  </div>
</template>

<script>
import TaskDetail from './components/TaskDetail'
export default {
  name: 'EditTask',
  components: { TaskDetail }
}
</script>

