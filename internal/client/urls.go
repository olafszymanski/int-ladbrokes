package client

const (
	classesUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/Class?translationLang=en&responseFormat=json&simpleFilter=class.isActive&simpleFilter=class.hasOpenEvent&simpleFilter=class.categoryId:equals:%v"
	eventsUrl  = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.isStarted:isFalse&simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"
)
