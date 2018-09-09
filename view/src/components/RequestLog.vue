<template>
  <div>
    <b-table striped hover :items="items"></b-table>
    <b-pagination size="sm" :total-rows="100" v-model="currentPage" :per-page="10">
    </b-pagination>
  </div>
</template>

<script>
  const items = [
    { isActive: true, age: 40, first_name: 'Dickerson', last_name: 'Macdonald' },
    { isActive: false, age: 21, first_name: 'Larsen', last_name: 'Shaw' },
    { isActive: false, age: 89, first_name: 'Geneva', last_name: 'Wilson' },
    { isActive: true, age: 38, first_name: 'Jami', last_name: 'Carney' }
  ]

  export default {
    data () {
      return {
        items: items,
        currentPage: 3
      }
    }
  }
</script>
