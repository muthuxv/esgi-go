package service

import broadcast "github.com/muthuxv/esgi-go/channels"

type Manager interface {
	OpenListener(paymentid string) chan interface{}
	CloseListener(paymentid string, channel chan interface{})
	Submit(userid, paymentid, text string)
	DeleteBroadcast(paymentid string)
}

type Message struct {
	UserId    string
	PaymentId string
	Text      string
}

type Listener struct {
	PaymentId string
	Chan      chan interface{}
}

type manager struct {
	paymentChannels map[string]broadcast.Broadcaster
	open            chan *Listener
	close           chan *Listener
	delete          chan string
	messages        chan *Message
}

var managerSingleton *manager

func GetPaymentManager() Manager {
	if managerSingleton == nil {
		managerSingleton = &manager{
			paymentChannels: make(map[string]broadcast.Broadcaster),
			open:            make(chan *Listener, 100),
			close:           make(chan *Listener, 100),
			delete:          make(chan string, 100),
			messages:        make(chan *Message, 100),
		}

		go managerSingleton.run()
	}

	return managerSingleton
}

func (m *manager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case paymentid := <-m.delete:
			m.deleteBroadcast(paymentid)
		case message := <-m.messages:
			m.payment(message.PaymentId).Submit(*message)
		}
	}
}

func (m *manager) register(listener *Listener) {
	m.payment(listener.PaymentId).Register(listener.Chan)
}

func (m *manager) deregister(listener *Listener) {
	m.payment(listener.PaymentId).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *manager) deleteBroadcast(paymentid string) {
	b, ok := m.paymentChannels[paymentid]
	if ok {
		b.Close()
		delete(m.paymentChannels, paymentid)
	}
}

/*
Get the payment with the id paymentid, or creates and registers it
*/
func (m *manager) payment(paymentid string) broadcast.Broadcaster {
	b, ok := m.paymentChannels[paymentid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.paymentChannels[paymentid] = b
	}
	return b
}

func (m *manager) OpenListener(paymentid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		PaymentId: paymentid,
		Chan:      listener,
	}
	return listener
}

func (m *manager) CloseListener(paymentid string, channel chan interface{}) {
	m.close <- &Listener{
		PaymentId: paymentid,
		Chan:      channel,
	}
}

func (m *manager) DeleteBroadcast(paymentid string) {
	m.delete <- paymentid
}

func (m *manager) Submit(userid, paymentid, text string) {
	msg := &Message{
		UserId:    userid,
		PaymentId: paymentid,
		Text:      text,
	}
	m.messages <- msg
}
