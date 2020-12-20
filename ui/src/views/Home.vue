<template>
  <div>
    <StoreCard
      name="元氣家族"
      priceLevel="c"
      description="元氣家族麵飯館"
      :categories="['麵', '飯']"
      :rating="3"
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
      :expanded="false"
    />
    <MenuCard name="白小白麵舖" :data="testMneu" />
  </div>
</template>

<script>
import StoreCard from "@/components/StoreCard";
import MenuCard from "@/components/MenuCard";

export default {
  name: "Home",
  components: {
    StoreCard,
    MenuCard,
  },
  data: () => ({
    testMneu: [
      {
        name: "經典豚骨拉麵",
        category: "拉麵",
        variants: [
          {
            name: "大碗",
            price: 148,
          },
          {
            name: "小碗",
            price: 118,
          },
        ],
      },
      {
        name: "味噌豚骨拉麵",
        category: "拉麵",
        variants: [
          {
            name: "大碗",
            price: 148,
          },
          {
            name: "小碗",
            price: 118,
          },
        ],
      },
      {
        name: "泡菜豚醬拌麵",
        category: "拌麵",
        variants: [
          {
            name: "大碗",
            price: 118,
          },
          {
            name: "小碗",
            price: 88,
          },
        ],
      },
      {
        name: "胡麻肉燥拌麵",
        category: "拌麵",
        variants: [
          {
            name: "大碗",
            price: 118,
          },
          {
            name: "小碗",
            price: 88,
          },
        ],
      },
      {
        name: "豆干",
        category: "單點",
        price: 20,
      },
      {
        name: "燙青菜",
        price: 30,
        category: "單點",
      },
      {
        name: "牛蒡天婦羅",
        price: 40,
        category: "炸物",
      },
    ],
  }),
};
</script>
