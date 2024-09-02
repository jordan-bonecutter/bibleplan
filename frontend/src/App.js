import {
    Button,
    Form,
    Input,
  DatePicker,
  Layout,
  message,
  Space,
  Typography
} from 'antd';
import { useState } from 'react';
import dayjs from 'dayjs';

const { Title, Paragraph, Text, Link } = Typography;

//const api = 'http://localhost:7771/api/v1/'
const api = 'https://bibleplan.jordanbonecutter.com/api/v1'

function App() {
  const [readings, setReadings] = useState([]);

  return <Layout style={{margin: 10}}>
    <Layout.Content>
    <Title>M'Cheyne Bible Reading Plan</Title>
    <Form onFinish={({email, date}) => {
      fetch(
        `${api}/plan`,
        {
          method: 'PUT',
          body: JSON.stringify({
            email,
            startDay: new Date(Date.parse(date)),
          }),
        },
      )
      .then(async resp => {
        if(resp.status === 200)
          message.success('Signed up for a daily bible plan')
        else
          message.error(`Failed siging up: ${await resp.text()}`)
      })
      .catch(err => message.error(`Failed siging up: ${err}`))
    }}>
      <Form.Item id="date" name="date" label="Date">
        <DatePicker minDate={dayjs().year(dayjs().year()-1)} maxDate={dayjs()} onChange={(e) => {
          fetch(`${api}/reading?start=${e.toDate().toJSON()}`)
            .then(r => r.json())
            .then(readings => setReadings(readings))
            .catch(console.log)
        }}/>
      </Form.Item>
      <Form.Item id="email" name="email" label="Email">
        <Space.Compact>
          <Input type="email"/>
          <Button type="primary" htmlType="submit">Submit</Button>
        </Space.Compact>
      </Form.Item>
    </Form>
    <>{readings.length === 0 ? <></> : <><br/><b>Today's Reading:</b><br/></>
    }</>
    <>{readings.map(reading => <><span>{reading}</span><br/></>)
    }</>
    <br/>
    <Paragraph>
    MY DEAR FLOCK , -- The approach of another year stirs up within me new desires for your
salvation, and for the growth of those of you who are saved. "God is my record how
greatly I long after you all in the bowels of Jesus Christ." What the coming year is to
bring forth who can tell? There is plainly a weight lying on the spirits of all good men,
and a looking for some strange work of judgment upon this land. There is a need now
to ask that solemn question -- "If in the land of peace wherein thou trustedst, they
wearied thee, then how wilt thou do in the swelling of Jordan?"</Paragraph>
    <Paragraph>
Those believers will stand firmest who have no dependence upon self or upon
creatures, but upon Jehovah our Righteousness. We must be driven more to our
Bibles, and to the mercy-seat, if we are to stand in the evil day. Then we shall be able
to say, like David – “The proud have had me greatly in derision, yet have I not declined
from Thy law.” “Princes have persecuted me without a cause, but my heart standeth in
awe of Thy word.”</Paragraph>
    <Paragraph>
It has long been in my mind to prepare a scheme of Scripture reading, in which as
many as were made willing by God might agree, so that the whole Bible might be read
once by you in the year, and all might be feeding in the same portion of the green
pasture at the same time.
I am quite aware that such a plan is accompanied with many
    </Paragraph>
    </Layout.Content>
  </Layout>
}

export default App;
