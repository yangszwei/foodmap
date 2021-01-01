<template>
  <v-card class="mx-auto my-12 pa-3" max-width="600">
    <v-card-title class="pb-0">{{ name }}</v-card-title>
    <v-card-text class="overline grey--text py-0 pb-1"
      >*菜單資訊僅供參考，實際以店家公告為主</v-card-text
    >
    <v-divider></v-divider>
    <div v-for="(category, index) in menuDataToCategories(data)" :key="index">
      <v-card-text class="py0 my-0 subtitle-2">
        {{ category }}
      </v-card-text>
      <v-card-text
        class="pt-0 pl-6 pb-3 grey--text text--darken-1 font-weight-medium"
        v-for="item in findCategoryItems(data, category)"
        :key="item.name"
      >
        <span>{{ item.name }}</span>
        <span class="float-right pr-3">
          <span v-for="price in itemToPrices(item)" :key="price" class="mx-2">
            {{ price }}
          </span>
        </span>
      </v-card-text>
    </div>
  </v-card>
</template>

<script>
import _ from "lodash";
export default {
  name: "MenuCard",
  props: {
    name: String,
    data: Array,
  },
  data: () => ({}),
  methods: {
    menuDataToCategories(__data) {
      return _.sortedUniq(_.values(_.mapValues(__data, "category")));
    },
    findCategoryItems(__data, __category) {
      return _.filter(__data, ["category", __category]);
    },
    itemToPrices(__item) {
      return __item.variants === undefined
        ? [__item.price]
        : _.reverse(
            _.sortedUniqBy(
              _.values(_.mapValues(__item.variants, "price")),
              parseFloat
            )
          );
    },
  },
};
</script>
