<template>
  <v-card class="mx-auto my-12" max-width="374">
    <v-card-title class="headline font-weight-medium">{{ name }}</v-card-title>
    <v-card-text>
      <v-row align="center" class="mx-0 my-0">
        <v-rating
          :value="rating"
          color="amber"
          dense
          background-color="grey lighten-1"
          half-increments
          size="1.5em"
        ></v-rating>

        <div class="grey--text ml-4">{{ rating }} ({{ comments.length }})</div>
      </v-row>

      <div class="my-4 subtitle-1">
        {{ priceLevelToDollarSigns(priceLevel) }} •
        <v-chip
          class="mx-1"
          small
          draggable
          v-for="categorie in categories"
          :key="categorie"
        >
          {{ categorie }}
        </v-chip>
      </div>

      <div v-if="description.length">{{ description }}</div>
    </v-card-text>
    <v-card-actions>
      <v-btn icon color="grey lighten-1">
        <v-icon>mdi-map-marker</v-icon>
      </v-btn>
      <v-btn icon color="grey lighten-1">
        <v-icon>mdi-clipboard-list</v-icon>
      </v-btn>
      <v-btn icon color="grey lighten-1">
        <v-icon>mdi-cards-heart</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <v-btn icon @click="show = !show">
        <v-icon>{{ show ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
      </v-btn>
    </v-card-actions>

    <v-expand-transition>
      <div v-if="show">
        <v-divider class="mx-4"></v-divider>
        <v-card-title>營業時間</v-card-title>
        <v-card-text>
          <v-chip-group
            v-model="selection"
            active-class="deep-purple accent-4 white--text"
            column
          >
            <v-chip>5:30PM</v-chip>
            <v-chip>7:30PM</v-chip>
            <v-chip>8:00PM</v-chip>
            <v-chip>9:00PM</v-chip>
          </v-chip-group>
        </v-card-text>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script>
export default {
  name: "StoreCard",
  props: {
    name: String,
    priceLevel: String,
    description: String,
    categories: Array,
    comments: Array,
    rating: Number,
  },
  data: () => ({
    show: false,
  }),
  methods: {
    priceLevelToDollarSigns(__priceLevel) {
      switch (__priceLevel) {
        case "c":
          return "$";
        case "m":
          return "$$";
        case "e":
          return "$$$";
      }
    },
    categoriesArrayToString(__categories) {
      return __categories.join(", ");
    },
  },
};
</script>
