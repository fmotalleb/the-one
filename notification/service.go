package notification

import (
	"context"

	"github.com/nikoksr/notify"
	"go.uber.org/zap"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/logging"
	"github.com/fmotalleb/the-one/notification/handlers"
)

var log = logging.LazyLogger("notification")

type Service struct {
	handlers map[string][]notify.Notifier
	bus      chan Notification
}

func New(cfg config.Config) (*Service, error) {
	log := log().Named("New")
	log.Info("Initializing Kernel")

	hb := map[string][]notify.Notifier{}
	for _, c := range cfg.Contacts {
		name := *c.Name.Unwrap()
		log.Debug("Processing contact", zap.String("name", name))

		services, err := handlers.FindHandler(c)
		if err != nil {
			log.Error("Failed to find handler", zap.String("name", name), zap.Error(err))
			return nil, err
		}

		if _, ok := hb[name]; !ok {
			hb[name] = make([]notify.Notifier, 0)
			log.Debug("Created new handler list", zap.String("name", name))
		}

		hb[name] = append(hb[name], services...)
		log.Debug("Added handlers", zap.String("name", name), zap.Int("count", len(services)))
	}

	service := new(Service)
	service.handlers = hb
	service.bus = make(chan Notification)
	log.Info("Kernel initialized successfully")
	go service.initWorker()
	return service, nil
}

func (k *Service) initWorker() {
	for n := range k.bus {
		go k.handleNotification(n)
	}
}

func (k *Service) handleNotification(n Notification) {
	log := log().Named("Handle")
	log.Debug("handling notification", zap.String("subject", n.Subject), zap.Strings("contacts", n.ContactPoints))

	for _, name := range n.ContactPoints {
		handlers, ok := k.handlers[name]
		if !ok {
			log.Warn("no handlers found for contact", zap.String("contact", name))
			continue
		}

		log.Debug("sending notification", zap.String("contact", name), zap.Int("handler_count", len(handlers)))

		notifier := notify.New()
		notifier.UseServices(handlers...)
		err := notifier.Send(
			n.Ctx,
			n.Subject,
			n.Message,
		)

		if err != nil {
			log.Error("failed to send notification", zap.String("contact", name), zap.Error(err))
		} else {
			log.Info("notification sent", zap.String("contact", name))
		}
	}
}

func (k *Service) Process(
	ctx context.Context,
	contactsPoints []string,
	subject, message string,
) {
	k.bus <- Notification{
		Ctx:           ctx,
		ContactPoints: contactsPoints,
		Subject:       subject,
		Message:       message,
	}
}
