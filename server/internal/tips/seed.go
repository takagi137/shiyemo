package tips

// SeedTips 种子数据：初始贴士列表
var SeedTips = []SurvivalTip{
	// 心态类
	{Content: "见过加班猝死的，没见过不上班饿死的。", Category: "mindset"},
	{Content: "今天不想投简历？那就不投。明天再说。", Category: "mindset"},
	{Content: "失业不是你的错，是结构性问题。", Category: "mindset"},
	{Content: "休息是为了更好地出发，不是偷懒。", Category: "mindset"},
	{Content: "别拿自己跟别人比，每个人的时区不同。", Category: "mindset"},
	{Content: "焦虑的时候，先把今天过好。", Category: "mindset"},
	{Content: "你的价值不由一份工作定义。", Category: "mindset"},
	{Content: "允许自己难过，但不要停下脚步。", Category: "mindset"},

	// 财务类
	{Content: "一个人生活自己做饭，每月 800 就能吃得很好。", Category: "finance"},
	{Content: "假如你有 10 万存款，每月花 3000，你能休息将近 3 年。", Category: "finance"},
	{Content: "失业金可以领 X 个月，去社保局申请就行。", Category: "finance"},
	{Content: "社保断缴可以自己续，医保只要连续缴满 X 个月就行。", Category: "finance"},
	{Content: "记账是控制支出的第一步。", Category: "finance"},
	{Content: "先砍掉不必要的订阅服务。", Category: "finance"},
	{Content: "二手平台卖掉闲置物品，既整理空间又有收入。", Category: "finance"},

	// 生活类
	{Content: "每天出门走走，哪怕只是去公园坐坐。", Category: "life"},
	{Content: "保持规律作息，这是你唯一能控制的事。", Category: "life"},
	{Content: "学一道新菜，给自己一点成就感。", Category: "life"},
	{Content: "图书馆是最好的免费空调办公区。", Category: "life"},
	{Content: "运动能缓解焦虑，跑步、游泳、散步都行。", Category: "life"},
	{Content: "整理房间，扔掉不需要的东西。", Category: "life"},
	{Content: "给自己设定一个小目标，比如每天读 30 页书。", Category: "life"},

	// 政策类
	{Content: "失业保险金申领条件：缴纳失业保险满 1 年，非因本人意愿中断就业。", Category: "policy"},
	{Content: "灵活就业人员可以自己缴纳社保和医保。", Category: "policy"},
	{Content: "部分地区有免费的职业技能培训，去人社局问问。", Category: "policy"},
	{Content: "创业担保贷款，最高可贷 20 万，政府贴息。", Category: "policy"},
	{Content: "公益性岗位是兜底保障，不要觉得丢人。", Category: "policy"},
}

// Seed 初始化种子数据
func Seed(repo TipsRepoInterface) error {
	for _, tip := range SeedTips {
		if err := repo.Create(&tip); err != nil {
			return err
		}
	}
	return nil
}
