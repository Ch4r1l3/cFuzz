<template>
  <div class="app-container">
    <task-create-edit :is-edit="true" />
  </div>
</template>

<script>
import TaskCreateEdit from './components/TaskCreateEdit'
export default {
  name: 'EditTask',
  components: { TaskCreateEdit }
}
</script>

