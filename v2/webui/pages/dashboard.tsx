import { gql, useQuery } from '@apollo/client';
import {
  Flex,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Text,
  Heading,
} from '@chakra-ui/react';
import { Layout } from '../components/Layout';
import { Navbar } from '../components/Navbar';

const tabList = ['Overview', 'Projections', 'Cluster', 'Security'];

const Dashboard = () => {
  const { data } = useQuery(
    gql`
      {
        uptime
        streamCount
        eventCount
      }
    `,
    {
      pollInterval: 1000,
    }
  );

  return (
    <Layout title="EventflowDB">
      <Flex h="100vh" bg="brand.500" flexDirection="column">
        <Navbar />

        <Flex bg="white" m={4} p={4} flexDirection="column">
          <Heading size="md" mb={4} color="brand.400">
            Dashboard
          </Heading>

          <Tabs>
            <TabList>
              {tabList.map((name, index) => (
                <Tab key={index} textTransform="uppercase" fontSize="sm">
                  {name}
                </Tab>
              ))}
            </TabList>

            <TabPanels>
              <TabPanel>
                <Flex>
                  <Text fontWeight="semibold" mr={2} color="brand.500">
                    Uptime:
                  </Text>
                  <Text>{data?.uptime}</Text>
                </Flex>

                <Flex>
                  <Text fontWeight="semibold" mr={2} color="brand.500">
                    Stream Count:
                  </Text>
                  <Text>{data?.streamCount}</Text>
                </Flex>

                <Flex>
                  <Text fontWeight="semibold" mr={2} color="brand.500">
                    Event Count:
                  </Text>
                  <Text>{data?.eventCount}</Text>
                </Flex>
              </TabPanel>
              <TabPanel>
                <p>two!</p>
              </TabPanel>
              <TabPanel>
                <p>three!</p>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Flex>
      </Flex>
    </Layout>
  );
};

export default Dashboard;
