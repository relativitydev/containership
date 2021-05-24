package azure

type ImagePromotion struct {
	SourceImage      string
	TargetRepository string
	SupportedTags    []string
	Destinations     []PromotionDestination
}

type PromotionDestination struct {
	Name           string
	Ring           int
	SubscriptionID string
	ResourceGroup  string
}
