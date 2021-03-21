import {
  Flex,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Heading,
} from '@chakra-ui/react';
import { ClusterTab } from '../components/ClusterTab';
import { Layout } from '../components/Layout';
import { Navbar } from '../components/Navbar';
import { OverviewTab } from '../components/OverviewTab';

const tabList = ['Overview', 'Projections', 'Cluster', 'Security'];

const Dashboard = () => {
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
                <OverviewTab />
              </TabPanel>

              <TabPanel>
                <p>two!</p>
              </TabPanel>

              <TabPanel>
                <ClusterTab />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Flex>
      </Flex>
    </Layout>
  );
};

export default Dashboard;
