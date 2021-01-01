<template>
  <v-card class="mx-auto my-12" max-width="400">
    <v-card-title class="headline font-weight-medium">{{ name }}</v-card-title>
    <v-card-text>
      <v-row align="center" class="mx-0 my-0">
        <v-rating
          :value="Math.round(rating)"
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
          v-for="category in categories"
          :key="category"
        >
          {{ category }}
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
      <v-btn icon @click="expanded = !expanded">
        <v-icon>{{ expanded ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
      </v-btn>
    </v-card-actions>

    <v-expand-transition>
      <div v-if="expanded">
        <v-divider class="mx-4"></v-divider>
        <v-card-title class="pb-0">營業時間</v-card-title>
        <v-card-title class="overline text--disabled py-0 pb-1">
          *營業時間可能因國定假日有所變化
        </v-card-title>
        <v-card-text>
          <div
            v-for="(business_hour, index) in business_hours"
            :key="index"
            class="subtitle-2 mx-2 my-1 font-weight-normal"
          >
            {{ businessHourIndexToString(index) }} •
            {{ businessHourToString(business_hour) }}
          </div>
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
    business_hours: Array,
    rating: Number,
    expanded: Boolean,
    menu: Object,
  },
  data: () => ({}),
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
    businessHourIndexToString(__business_hour_index) {
      const weekday = [
        "星期一",
        "星期二",
        "星期三",
        "星期四",
        "星期五",
        "星期六",
        "星期日",
      ];
      return weekday[__business_hour_index];
    },
    __timeParser(__timeString) {
      const __time = __timeString.split(":");
      const hour =
        __time[0] < 12
          ? `上午 ${__time[0].lpad("0", 2)}`
          : `下午 ${String(__time[0] - 12).lpad("0", 2)}`;
      const minute = __time[1];
      return `${hour}:${minute}`;
    },
    businessHourToString(__business_hour) {
      let openTime = this.__timeParser(__business_hour[0]);
      let closeTime = this.__timeParser(__business_hour[1]);
      return `${openTime} ~ ${closeTime}`;
    },
  },
};
</script>
