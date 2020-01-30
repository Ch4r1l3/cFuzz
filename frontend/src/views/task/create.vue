<template>
  <div class="app-container">
    <task-detail />
  </div>
</template>

<script>
import TaskDetail from './components/TaskDetail'
export default {
  name: 'CreateTask',
  components: { TaskDetail }
}
</script>

