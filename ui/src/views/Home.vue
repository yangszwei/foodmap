<template>
  <StoreCard
    name="元氣家族"
    priceLevel="c"
    description="元氣家族麵飯館"
    :categories="['麵', '飯']"
    rating="3"
    :business_hours="[
      ['07:00', '20:00'],
      ['07:00', '20:00'],
      ['07:00', '20:00'],
      ['07:00', '20:00'],
      ['07:00', '20:00'],
      ['07:00', '20:00'],
      ['07:00', '20:00'],
    ]"
    :comments="[]"
    :menu="{}"
    :expanded="true"
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
