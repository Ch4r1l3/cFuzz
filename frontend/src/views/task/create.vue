<template>
  <div class="app-container">
    <task-create-edit />
  </div>
</template>

<script>
import TaskCreateEdit from './components/TaskCreateEdit'
export default {
  name: 'CreateTask',
  components: { TaskCreateEdit }
}
</script>

