<template>
  <StoreCard
    name="元氣家族"
    priceLevel="c"
    description="元氣家族麵飯館"
    :categories="['麵', '飯']"
    rating="3"
    :comments="[]"
  />
</template>

<script>
import StoreCard from "@/components/StoreCard";

export default {
  name: "Home",
  components: {
    StoreCard,
  },
};
</script>
