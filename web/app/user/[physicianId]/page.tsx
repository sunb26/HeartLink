import { type Patient, columns } from "@/components/ui/patients/columns";
import { DataTable } from "@/components/ui/patients/data-table";
import { currentUser } from "@clerk/nextjs/server";

async function getData(): Promise<Patient[]> {
  // Fetch data from your API here.
  return [
      {
        id: 1,
        firstName: "John",
        lastName: "Smith",
        email: "john.smith@example.com",
        lastUpdated: "2024-02-15T09:30:00Z",
        verified: true,
      },
      {
        id: 2,
        firstName: "Sarah",
        lastName: "Johnson",
        email: "sarah.j@example.com",
        lastUpdated: "2024-02-14T15:45:00Z",
        verified: true,
      },
      {
        id: 3,
        firstName: "Miguel",
        lastName: "Rodriguez",
        email: "mrodriguez@example.com",
        lastUpdated: "2024-02-13T11:20:00Z",
        verified: true,
      },
      {
        id: 4,
        firstName: "Emma",
        lastName: "Wilson",
        email: "emma.w@example.com",
        lastUpdated: "2024-02-12T16:15:00Z",
        verified: true,
      },
      {
        id: 5,
        firstName: "Chen",
        lastName: "Wei",
        email: "chen.wei@example.com",
        lastUpdated: "2024-02-11T08:45:00Z",
        verified: true,
      },
      {
        id: 6,
        firstName: "Priya",
        lastName: "Patel",
        email: "priya.p@example.com",
        lastUpdated: "2024-02-10T14:30:00Z",
        verified: true,
      },
      {
        id: 7,
        firstName: "Alex",
        lastName: "Thompson",
        email: "alex.t@example.com",
        lastUpdated: "2024-02-09T10:20:00Z",
        verified: true,
      },
      {
        id: 8,
        firstName: "Maria",
        lastName: "Garcia",
        email: "mgarcia@example.com",
        lastUpdated: "2024-02-08T13:15:00Z",
        verified: false,
      },
      {
        id: 9,
        firstName: "James",
        lastName: "Wilson",
        email: "jwilson@example.com",
        lastUpdated: "2024-02-07T17:40:00Z",
        verified: false,
      },
      {
        id: 10,
        firstName: "Lisa",
        lastName: "Anderson",
        email: "lisa.a@example.com",
        lastUpdated: "2024-02-06T12:10:00Z",
        verified: false,
      },
      {
        id: 11,
        firstName: "David",
        lastName: "Kim",
        email: "dkim@example.com",
        lastUpdated: "2024-02-05T09:55:00Z",
        verified: false,
      }
    // ...
  ];
}

export default async function PhysicianPage() {
  const physician = await currentUser();
  const data = await getData();

  return (
    <div>
      <div className="bg-off-white">
        <h1 className="text-4xl font-bold text-center py-10 font-[Syne]">
          Welcome {physician?.fullName}
        </h1>
      </div>
      <div className="container mx-auto py-10">
        <h2 className="text-2xl font-bold pb-4 font-[Syne]">Your Patients</h2>
        <DataTable columns={columns} data={data} />
      </div>
    </div>
  );
}
